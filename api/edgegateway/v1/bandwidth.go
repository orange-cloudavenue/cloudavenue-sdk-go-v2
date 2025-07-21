package edgegateway

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/common-go/validators"
)

type ( //nolint:gocritic
	classOfService struct {
		// MaxBandwidth defines the maximum bandwidth in Mbps for this class of service.
		MaxBandwidth int

		// MaxEdgeGateways defines the maximum number of edge gateways allowed for this class of service.
		MaxEdgeGateways int

		// MaxEdgeGatewayBandwidth defines the values of bandwidth that can be allocated to each edge gateway.
		// This is a list of integers representing the allowed bandwidth values in Mbps.
		MaxEdgeGatewayBandwidth []int

		// Allow unlimited bandwidth for edge gateways.
		AllowUnlimited bool
	}
)

// Source : https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/network/
var classOfServices = map[string]classOfService{
	"SHARED_STANDARD": {
		MaxBandwidth:            300,
		MaxEdgeGateways:         4,
		MaxEdgeGatewayBandwidth: []int{5, 25, 50, 75, 100, 150, 200, 250, 300},
		AllowUnlimited:          false,
	},
	"SHARED_PREMIUM": {
		MaxBandwidth:            1000,
		MaxEdgeGateways:         8,
		MaxEdgeGatewayBandwidth: []int{5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000},
		AllowUnlimited:          false,
	},
	"DEDICATED_MEDIUM": {
		MaxBandwidth:            3500,
		MaxEdgeGateways:         6,
		MaxEdgeGatewayBandwidth: []int{5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000, 2000},
		AllowUnlimited:          true, // Unlimited bandwidth is allowed for edge gateway with this class of service.
	},
	"DEDICATED_LARGE": {
		MaxBandwidth:            10000,
		MaxEdgeGateways:         12,
		MaxEdgeGatewayBandwidth: []int{5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 6000},
		AllowUnlimited:          true, // Unlimited bandwidth is allowed for edge gateway with this class of service.
	},
}

func (c *Client) GetBandwidth(ctx context.Context, params ParamsEdgeGateway) (*ModelBandwidth, error) {
	logger := c.logger.WithGroup("GetBandwidth")

	if err := validators.New().Struct(&params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// ID is required to request the API.
	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return nil, err
		}
	}

	// Get the endpoint for the edge gateway
	// Error is ignored here because the endpoint is registered at package init time.
	ep, _ := cav.GetEndpoint("Bandwidth", cav.MethodGET)

	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.ID))
	if err != nil {
		logger.Error("Failed to get edge gateway bandwidth", "error", err)
		return nil, err
	}

	return resp.Result().(*apiResponseBandwidth).toModel(), nil
}

func (c *Client) GetEdgeGatewayBandwidth(ctx context.Context, params ParamsEdgeGateway) (*ModelT0EdgeGateway, error) {
	logger := c.logger.WithGroup("GetEdgeGatewayBandwidth")

	// GetT0FromEdgeGateway is called to get T0
	t0, err := c.GetT0FromEdgeGateway(ctx, params)
	if err != nil {
		logger.Error("Failed to get T0 from edge gateway", "error", err)
		return nil, err
	}

	var edgeGateway *ModelT0EdgeGateway
	for _, eg := range t0.EdgeGateways {
		if eg.ID == params.ID || eg.Name == params.Name {
			edgeGateway = &eg
			break
		}
	}

	// edgeGateway is never nil because GetT0FromEdgeGateway ensures the edge gateway exists

	return edgeGateway, nil
}
