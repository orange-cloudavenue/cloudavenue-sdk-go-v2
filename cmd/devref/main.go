package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

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

func main() {
	// Get configuration from environment
	baseURL := os.Getenv("DEVFLOW_API_URL")
	if baseURL == "" {
		baseURL = "https://devflow.cav007.myaddr.tools"
	}

	if len(os.Args) != 2 {
		log.Fatal("code argument is required (e.g., '123')")
	}

	inputArg := os.Args[1]
	if inputArg == "" {
		log.Fatal("code argument is required (e.g., '123')")
	}

	var branch string

	// Create the client with OAuth Device Flow authentication
	client, err := sdk.NewClient(
		baseURL,
		sdk.WithInsecure(),
		sdk.WithAutoRefresh(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Authenticate if not already authenticated
	if !client.IsAuthenticated() {
		fmt.Println("No authentication token found. Starting device flow...")
		err := client.Authenticate(context.Background(), "devref-v2", func(userCode, verificationURI string) {
			fmt.Printf("\n🔐 Authentication required\n")
			fmt.Printf("   Visit: %s\n", verificationURI)
			fmt.Printf("   Enter code: %s\n", userCode)
			fmt.Println("   Waiting for authorization...")
		})
		if err != nil {
			log.Fatal("Authentication failed:", err)
		}
		fmt.Println("✅ Authentication successful!")
	}

	ctx := context.Background()

	var code int
	fmt.Sscanf(inputArg, "%d", &code)
	resourceInfo, err := client.GetResourceByCode(ctx, code)
	if err != nil {
		log.Fatalf("❌ Error resolving code %03d: %v", code, err)
	}
	fmt.Printf("✅ Code %03d resolved to: %s/%s on branch %s\n", code, resourceInfo.NamespaceName, resourceInfo.ResourceName, resourceInfo.Branch)
	namespaceName := resourceInfo.NamespaceName
	resourceName := resourceInfo.ResourceName
	branch = resourceInfo.Branch

	// // Allocate a unique code for this namespace/resource/branch combination
	// code, err = client.AllocateCode(ctx, namespaceID, resourceName, branch)
	// if err != nil {
	// 	log.Fatal("❌ Error allocating code:", err)
	// }

	// if code != nil {
	// 	fmt.Printf("✅ Code alloué : %03d pour %s/%s\n", *code, resourceName, branch)
	// } else {
	// 	fmt.Println("✅ Aucun code alloué (branche main)")
	// }

	x := reg.GetCommandsByFilter(func(cmd commands.Command) bool {
		if namespaceName == resourceName {
			return cmd.GetNamespace() == namespaceName && cmd.GetResource() == "" && cmd.GetVerb() != ""
		}

		return cmd.GetNamespace() == namespaceName && cmd.GetResource() == resourceName && cmd.GetVerb() != ""
	})

	if len(x) == 0 {
		log.Fatalf("No commands found for namespace '%s' and resource '%s'", namespaceName, resourceName)
	}

	// Define the SDK configuration with required resource identifiers
	config := &models.SDKConfig{
		ResourceID:  resourceInfo.ResourceID,
		NamespaceID: resourceInfo.NamespaceID,
		Branch:      branch,
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

	resp, err := client.Upload(ctx, config)
	if err != nil {
		log.Fatal("Error uploading SDK config:", err)
	}

	// Display results including code
	fmt.Println("\n✅ SDK upload réussi !")
	fmt.Printf("   Resource : %s\n", resp.Resource.Name)
	fmt.Printf("   Branch   : %s\n", resp.Resource.Branch)
	if resp.Resource.Code != nil {
		fmt.Printf("   Code     : %03d\n", *resp.Resource.Code)
	}
}

func convertParams(params pspecs.Params) []models.SDKAttribute {
	if params == nil {
		return nil
	}

	return convertParamsSpecs(params)
}

func convertParamsSpecs(paramsSpecs pspecs.Params) []models.SDKAttribute {
	if paramsSpecs == nil {
		return nil
	}

	var attrs []models.SDKAttribute

	for _, p := range paramsSpecs {
		attr := models.SDKAttribute{
			Name:        p.GetName(),
			Description: p.GetDescription(),
			Type:        getTypeString(p.GetType().Type()),
			Validations: make([]models.SDKValidation, 0),
		}

		// Handle nested structs, slices, and maps
		if nested, ok := p.(pspecs.ParamSpecNested); ok {
			children := convertParamsSpecs(nested.GetItemsSpec())
			if len(children) > 0 {
				attr.Children = children
			}
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
