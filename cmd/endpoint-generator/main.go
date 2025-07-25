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
	"strings"
	"text/template"
)

type endpoints []endpoint
type endpoint struct {
	Package          string
	DocumentationURL string
	Description      string
	Name             string
	SubClient        string
	PathTemplate     string
	BodyRequestType  string
	BodyResponseType string
}

//go:embed generator.tmpl
var tmplFile embed.FS

func main() {

	var (
		flagPath     = flag.String("path", "", "The path to the file to generate commands from")
		flagFilename = flag.String("filename", "", "The name of the file to generate")
		flagDebug    = flag.Bool("debug", false, "Enable debug mode")
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

	endpts := make(endpoints, 0)

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

				// Check if the call is to Register
				fun, ok := c.Fun.(*ast.SelectorExpr)
				if !ok || fun.Sel.Name != "Register" {
					continue
				}

				xFun, ok := fun.X.(*ast.CompositeLit)
				if !ok {
					continue
				}

				endpoint := &endpoint{
					Package: fileAst.Name.Name,
				}

				endpValue := reflect.ValueOf(endpoint).Elem()

				for _, arg := range xFun.Elts {
					kv, ok := arg.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					key, ok := kv.Key.(*ast.Ident)
					if !ok {
						continue
					}

					fieldValue := endpValue.FieldByName(key.Name)
					if fieldValue.IsValid() && fieldValue.CanSet() {
						fieldValue.SetString(findValue(kv))
					}
				}

				endpts = append(endpts, *endpoint)
			}

		}
	}

	var endpointTmpl = struct {
		PackageName string
		Endpoints   endpoints
	}{
		PackageName: fileAst.Name.Name,
		Endpoints:   endpts,
	}

	tmpl, err := template.ParseFS(tmplFile, "generator.tmpl")
	if err != nil {
		panic(err)
	}

	var outputPath string

	split := strings.Split(pwd+"/"+*flagPath, "/")
	findOutputDir := func() string {
		// src: /github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1/
		// to: /github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/
		for i := len(split) - 1; i >= 0; i-- {
			if split[i] == "cloudavenue-sdk-go-v2" {
				return strings.Join(split[:i+1], "/") + "/endpoints/"
			}
		}
		return ""
	}

	if *flagFilename == "" {

		// flagPath == /api/edgegateway/v1/edgegateway_endpoints.go
		// We want to generate the file in /api/edgegateway/v1/zz_<namespace>.go
		nameExtracted := split[len(split)-1]
		nameExtracted = nameExtracted[:len(nameExtracted)-len("_endpoints.go")]

		outputPath = findOutputDir() + "zz_" + nameExtracted + ".go"

	} else {
		outputPath = findOutputDir() + *flagFilename
	}

	if *flagDebug {
		log.Default().Print("Path to output file: ", outputPath)
	}

	// Create io.Writer to write the output to file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, endpointTmpl)
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
	case *ast.CompositeLit:
		kvc, ok := v.Type.(*ast.Ident)
		if !ok {
			return ""
		}
		return kvc.Name
	default:
		return ""
	}
}
