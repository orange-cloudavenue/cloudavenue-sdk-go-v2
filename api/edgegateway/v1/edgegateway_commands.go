package edgegateway

import (
	"context"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate command-generator -path edgegateway_commands.go

func init() {
	// * GetEdgeGateway
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Verb:      "Get",

		ShortDocumentation: "GetEdgeGateway retrieves an edge gateway",
		LongDocumentation:  "Get EdgeGateway performs a GET request to retrieve an edge gateway",

		ParamsType: ParamsEdgeGateway{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.RequireIfParamIsNull("Name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgeGateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.RequireIfParamIsNull("ID"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		ModelType: ModelEdgeGateway{},

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsEdgeGateway)
			ep, _ := cav.GetEndpoint("EdgeGateway", cav.MethodGET)

			logger := cc.logger.WithGroup("GetEdgeGateway")

			// ID is required to request the API.
			if p.ID == "" {
				var err error
				p.ID, err = cc.retrieveEdgeGatewayIDByName(ctx, p.Name)
				if err != nil {
					return nil, err
				}
			}

			resp, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.ID),
			)
			if err != nil {
				logger.Error("Failed to get edge gateway", "error", err)
				return nil, err
			}

			return resp.Result().(*apiResponseEdgegateway).toModel(), nil
		},
		AutoGenerate: true,
	})

	// * ListEdgeGateway
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Verb:               "List",
		ShortDocumentation: "ListEdgeGateways retrieves a list of edge gateways",
		LongDocumentation:  "List EdgeGateways performs a GET request to retrieve a list of edge gateways",
		ModelType:          ModelEdgeGateways{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			ep, _ := cav.GetEndpoint("ListEdgeGateway", cav.MethodGET)

			logger := cc.logger.WithGroup("ListEdgeGateways")

			resp, err := cc.c.Do(
				ctx,
				ep,
			)
			if err != nil {
				logger.Error("Failed to list edge gateways", "error", err)
				return nil, err
			}

			return resp.Result().(*apiResponseEdgegateways).toModel(), nil
		},
		AutoGenerate: true,
	})

	// * CreateEdgeGateway
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Verb:               "Create",
		ShortDocumentation: "CreateEdgeGateway creates a new edge gateway",
		LongDocumentation:  "Create EdgeGateway performs a POST request to create a new edge gateway",
		AutoGenerate:       true,
		ParamsType:         ParamsCreateEdgeGateway{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "ownerType",
				Description: "The type of the owner of the edge gateway.",
				Required:    true,
				Validators: []commands.Validator{
					commands.ValidatorOneOf("vdc", "vdcgroup"),
				},
			},
			commands.ParamsSpec{
				Name:        "ownerName",
				Description: "The name of the VDC or VDC Group that this edge gateway belongs to.",
				Required:    true,
			},
			commands.ParamsSpec{
				Name:        "t0Name",
				Description: "The name of the T0 router that this edge gateway will be connected to.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					// TODO validator cav name
				},
			},
		},
		ModelType: ModelEdgeGateway{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsCreateEdgeGateway)
			ep, _ := cav.GetEndpoint("EdgeGateway", cav.MethodPOST)

			logger := cc.logger.WithGroup("CreateEdgeGateway")

			if p.T0Name == "" {
				// If T0Name is not provided, retrieve the first available T0 router.
				t0s, err := cc.ListT0(ctx)
				if err != nil {
					return nil, err
				}
				if t0s.Count == 0 {
					logger.Error("No T0 routers available to connect the edge gateway")
					return nil, errors.New("No T0 routers available to connect the edge gateway")
				}
				if t0s.Count > 1 {
					logger.Warn("Multiple T0 routers found, using the first one", "count", t0s.Count)
				}
				p.T0Name = t0s.T0s[0].Name
			}

			// Prepare the request body
			reqBody := apiRequestEdgeGateway{
				T0Name: p.T0Name,
			}

			var edgeGatewayCreated string

			ep.SetJobExtractorFunc(cav.ExtractorFunc(func(resp *resty.Response) {
				r, ok := resp.Result().(*cav.CerberusJobAPIResponse)
				if !ok {
					logger.Error("Failed to extract job information")
					return
				}

				if len(*r) == 0 {
					logger.Error("No job information returned")
					return
				}

				job := (*r)[0]
				for _, j := range job.Actions {
					if err := validators.New().Var(j.Details, "edgegateway_name"); err == nil {
						edgeGatewayCreated = j.Details
						break
					}
				}
			}))

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.OwnerType),
				cav.WithPathParam(ep.PathParams[1], p.OwnerName),
				cav.SetBody(reqBody),
			)
			if err != nil {
				logger.Error("Failed to create edge gateway", "error", err)
				return nil, err
			}

			return cc.GetEdgeGateway(ctx, ParamsEdgeGateway{
				Name: edgeGatewayCreated,
			})
		},
	})
}
