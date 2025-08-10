package edgegateway

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
)

//go:generate command-generator -path bandwidth_commands.go

func init() {
	// Register the GetEdgeGatewayBandwidth command
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Resource:           "Bandwidth",
		Verb:               "Get",
		ShortDocumentation: "Get the bandwidth of an edge gateway.",
		LongDocumentation:  "Get the bandwidth of an edge gateway. This command retrieves the bandwidth information for a specific edge gateway.",
		AutoGenerate:       true,

		ParamsType: ParamsEdgeGateway{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("Name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("ID"),
					commands.ValidatorOmitempty(),
				},
			},
		},

		ModelType: ModelEdgeGatewayBandwidth{},

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsEdgeGateway)

			pT0 := ParamsGetT0{
				EdgegatewayID:   p.ID,
				EdgegatewayName: p.Name,
			}

			t0, err := cc.GetT0(ctx, pT0)
			if err != nil {
				return nil, err
			}

			var edgeGateway *ModelT0EdgeGateway
			for _, eg := range t0.EdgeGateways {
				if eg.ID == p.ID || eg.Name == p.Name {
					edgeGateway = &eg
					break
				}
			}

			// edgeGateway is never nil because GetT0FromEdgeGateway ensures the edge gateway exists

			// Transform ModelT0EdgeGateway to ModelEdgeGatewayBandwidth
			// Create a new ModelEdgeGatewayBandwidth from ModelT0EdgeGateway
			bandwidth := &ModelEdgeGatewayBandwidth{
				ID:                     edgeGateway.ID,
				Name:                   edgeGateway.Name,
				Bandwidth:              edgeGateway.Bandwidth,
				AllowedBandwidthValues: edgeGateway.AllowedBandwidthValues,
			}

			return bandwidth, nil
		},
	})
}
