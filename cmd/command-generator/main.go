package main

import (
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/kr/pretty"
)

type cmds []*cmd

type cmd struct {
	Package                   string
	Namespace, Resource, Verb string
	ParamsType, ModelType     string
	AutoGenerate              bool
}

//go:embed generator.tmpl
var tmplFile embed.FS

// TODO match AutoGenerate: true in the command definition

func main() {

	var (
		flagPath  = flag.String("path", "", "The path to the file to generate commands from")
		flagDebug = flag.Bool("debug", false, "Enable debug mode")
	)

	flag.Parse()

	if *flagPath == "" {
		flag.Usage()
		return
	}

	// print pwd
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	commands := make(cmds, 0)

	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, pwd+"/"+*flagPath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	if *flagDebug {
		ast.Print(fset, fileAst)
	}

	for _, decl := range fileAst.Decls {
		f, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if the function is init
		if f.Name.Name == "init" {
			// Generate commands for the init function

			for _, stmt := range f.Body.List {
				e, ok := stmt.(*ast.ExprStmt)
				if !ok {
					continue
				}

				c, ok := e.X.(*ast.CallExpr)
				if !ok {
					continue
				}

				// Check if the call is to Cmds.Register
				fun, ok := c.Fun.(*ast.SelectorExpr)
				if !ok || fun.Sel.Name != "Register" {
					continue
				}

				// Check if the first argument is a composite literal
				if len(c.Args) == 1 {

					cmd := &cmd{
						Package: fileAst.Name.Name,
					}

					for _, arg := range c.Args {
						compLit, ok := arg.(*ast.CompositeLit)
						if !ok {
							continue
						}

						cmdValue := reflect.ValueOf(cmd).Elem()

						// Retrieve the Namespace, Resource and Verb from the fields
						for _, field := range compLit.Elts {
							kv, ok := field.(*ast.KeyValueExpr)
							if !ok {
								continue
							}

							key, ok := kv.Key.(*ast.Ident)
							if !ok {
								continue
							}

							fieldValue := cmdValue.FieldByName(key.Name)
							if fieldValue.IsValid() && fieldValue.CanSet() {
								v := findValue(kv)
								value := reflect.ValueOf(v)
								switch fieldValue.Kind() {
								case reflect.String:
									if value.Kind() == reflect.String {
										fieldValue.SetString(value.String())
									}
								case reflect.Bool:
									boolValue := false
									switch value.Kind() {
									case reflect.Bool:
										boolValue = value.Bool()
									case reflect.String:
										// convert string to bool
										bv, err := strconv.ParseBool(v)
										if err != nil {
											continue
										}

										boolValue = bv
									}

									fieldValue.SetBool(boolValue)
								case reflect.Int, reflect.Int64:
									if value.Kind() == reflect.Int || value.Kind() == reflect.Int64 {
										intValue := value.Int()
										fieldValue.SetInt(intValue)
									}
								}
							}

						}
					}

					if cmd.AutoGenerate {
						commands = append(commands, cmd)
					}
				}
			}

		}
	}

	if *flagDebug {
		pretty.Print("Commands", commands)
	}

	var commandTmpl = struct {
		PackageName string
		Commands    cmds
	}{
		PackageName: fileAst.Name.Name,
		Commands:    commands,
	}

	tmpl, err := template.ParseFS(tmplFile, "generator.tmpl")
	if err != nil {
		panic(err)
	}

	// flagPath == /api/edgegateway/v1/edgegateway_commands.go
	// We want to generate the file in /api/edgegateway/v1/zz_<namespace>.go
	split := strings.Split(pwd+"/"+*flagPath, "/")
	nameExtracted := split[len(split)-1]
	nameExtracted = nameExtracted[:len(nameExtracted)-len("_commands.go")]

	outputPath := strings.Join(split[:len(split)-1], "/") + "/zz_" + nameExtracted + ".go"

	log.Default().Print("Path to output file: ", outputPath)

	// Create io.Writer to write the output to file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, commandTmpl)
	if err != nil {
		panic(err)
	}

}

func cleanQuote(s string) string {
	if len(s) < 2 {
		return s
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func findValue(kv *ast.KeyValueExpr) string {
	if kv == nil {
		return ""
	}

	switch v := kv.Value.(type) {
	case *ast.BasicLit:
		return cleanQuote(v.Value)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", v.X, v.Sel.Name)
	case *ast.Ident:
		return cleanQuote(v.Name)
	case *ast.CompositeLit:
		kvc, ok := v.Type.(*ast.Ident)
		if !ok {
			fmt.Println("Could not find type for composite literal:", v.Type)
			return ""
		}
		return cleanQuote(kvc.Name)
	default:
		return ""
	}
}
