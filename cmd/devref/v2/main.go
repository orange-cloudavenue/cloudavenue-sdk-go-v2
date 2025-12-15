package main

import (
	"context"
	"log"
	"reflect"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/pspecs"
	"github.com/orange-cloudavenue/devflow-sdk-go/models"
	"github.com/orange-cloudavenue/devflow-sdk-go/sdk"

	// Force import of all commands to register them
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/draas/v1"
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/organization/v1"
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdcgroup/v1"
)

// NewRegistry get the global command registry
var reg = commands.NewRegistry()

var (
	argNamespace = "EdgeGateway"
	argResource  = ""
)

func main() {
	// Create the client
	client, err := sdk.NewClient(
		"http://localhost:3000",
		sdk.WithToken("sdk_InlxyTDvwhnIrEOSCyIeNzobBQBBxF3h"),
		sdk.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	x := reg.GetCommandsByFilter(func(cmd commands.Command) bool {
		return cmd.GetNamespace() == argNamespace && cmd.GetResource() == argResource && cmd.GetVerb() != ""
	})

	if len(x) == 0 {
		log.Fatalf("No commands found for namespace '%s' and resource '%s'", argNamespace, argResource)
	}

	// Define the SDK configuration
	config := &models.SDKConfig{
		Language:    models.LanguageGo,
		PackageName: "cloudavenue",
	}

	for _, cmd := range x {
		config.Commands = append(config.Commands, models.SDKCommandMapping{
			Verb:          cmd.GetVerb(),
			SDKMethodName: cmd.AutoGenerateCustomFuncName,
			Params:        convertParams(cmd.ParamsSpecs),
			Response:      convertModel(cmd.ModelType),
		})
	}

	ctx := context.Background()

	pp.Print(config)

	resp, err := client.Upload(ctx, config)
	if err != nil {
		log.Fatal("Error uploading SDK config:", err)
	}

	log.Println("SDK upload response:", resp)
}

func convertParams(params pspecs.Params) []models.SDKAttribute {
	var attrs []models.SDKAttribute

	for _, p := range params {
		attr := models.SDKAttribute{
			Name:        p.GetName(),
			Description: p.GetDescription(),
			Type:        p.GetType().Type().String(),
			Validations: make([]models.SDKValidation, 0),
		}

		// Special case required
		if p.IsRequired() {
			attr.Validations = append(attr.Validations, models.SDKValidation{
				Type: models.ValidationRequired,
			})
		}

		// Add validations
		for _, v := range p.GetValidators() {
			switch v.(type) {
			default:
				attr.Validations = append(attr.Validations, models.SDKValidation{
					Type:    models.ValidationCustom,
					Message: v.GetMarkdownDescription(),
					Value:   v.GetKey(),
				})
			}
		}

		attrs = append(attrs, attr)
	}

	return attrs
}

func convertModel(modelType any) []models.SDKAttribute {
	if modelType == nil {
		return nil
	}

	t := reflect.TypeOf(modelType)
	// Dereference pointer
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	return walkStruct(t)
}

func walkStruct(t reflect.Type) []models.SDKAttribute {
	var attrs []models.SDKAttribute

	// Dereference pointer
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get json tag or fallback to field name
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = field.Name
		} else {
			// Remove omitempty or other tag options
			if idx := strings.Index(fieldName, ","); idx != -1 {
				fieldName = fieldName[:idx]
			}
		}

		// Skip fields with json:"-"
		if fieldName == "-" {
			continue
		}

		// Handle anonymous (embedded) fields - merge their fields into parent
		if field.Anonymous {
			embeddedAttrs := walkType(field.Type)
			attrs = append(attrs, embeddedAttrs...)
			continue
		}

		// Get documentation from tag
		description := field.Tag.Get("documentation")

		attr := models.SDKAttribute{
			Name:        fieldName,
			Description: description,
			Type:        getTypeString(field.Type),
			Validations: make([]models.SDKValidation, 0),
		}

		// Handle nested structs, slices, and maps
		children := walkType(field.Type)
		if len(children) > 0 {
			attr.Children = children
		}

		attrs = append(attrs, attr)
	}

	return attrs
}

func walkType(t reflect.Type) []models.SDKAttribute {
	// Dereference pointer
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		return walkStruct(t)
	case reflect.Slice, reflect.Array:
		elemType := t.Elem()
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() == reflect.Struct {
			return walkStruct(elemType)
		}
	case reflect.Map:
		elemType := t.Elem()
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() == reflect.Struct {
			return walkStruct(elemType)
		}
	}

	return nil
}

func getTypeString(t reflect.Type) string {
	// Dereference pointer
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:
		return "bool"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map:
		return "object"
	case reflect.Struct:
		return "object"
	default:
		return t.String()
	}
}
