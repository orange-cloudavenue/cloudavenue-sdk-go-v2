package commands

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

type Command struct {
	// Namespace is the namespace of the command, e.g. "edgegateway", "vdc", "vapp", etc.
	// It is used to group commands together and to define the endpoint for the command.
	// If no resource is defined, the namespace is used as the resource.
	Namespace string `validate:"required"`

	// Resource is the resource of the command, e.g. "firewall-rule, "security-group", etc.
	Resource string `validate:"omitempty"`

	// Verb is the verb of the command, e.g. "create", "delete", "update", "list", "get".
	Verb string `validate:"required"`

	// Documentation is the documentation of the command.
	ShortDocumentation string `validate:"required"`
	LongDocumentation  string `validate:"omitempty"`

	// Paramspec defines specifications for arguments.
	ParamsSpecs ParamsSpecs

	// ParamsType defines the type of parameters for the command.
	ParamsType any

	// ModelType defines the type of the model returned by the command.
	ModelType any

	// Runner
	RunnerFunc func(ctx context.Context, cmd *Command, client, params any) (any, error)

	// Deprecated defines whether the command is deprecated.
	Deprecated        bool
	DeprecatedMessage string

	// AutoGenerate defines whether the command is auto-generated.
	AutoGenerate bool
}

// Func
type (
	RunnerFunc          func(ctx context.Context, ep *cav.Endpoint, client, params any) (any, error)
	ParamsValidatorFunc func(ctx context.Context, value any, paramsSpecs ParamsSpecs, params any) error
	ParamsTransformFunc func(value string, paramsSpecs ParamsSpecs, params any) (string, error)
)

// * Parameters

type (
	ParamsSpecs []ParamsSpec

	ParamsSpec struct {
		// Name is the name of the argument.
		Name string `validate:"required"`

		// Description is the description of the argument.
		Description string `validate:"required"`

		// Required defines whether the argument is required.
		Required bool

		// Validator is a function that validates the argument value.
		Validators []Validator

		// TransformFunc is a function that transforms the argument value.
		// It is used to transform the value after validation.
		TransformFunc ParamsTransformFunc
	}
)

// // GetParamsSpec returns the ParamsSpec for the given name.
// func (ps *ParamsSpecs) GetParamsSpec(name string) *ParamsSpec {
// 	for _, spec := range *ps {
// 		if spec.Name == name {
// 			return &spec
// 		}
// 	}
// 	return nil
// }
