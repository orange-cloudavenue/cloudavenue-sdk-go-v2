package edgegateway

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/common-go/validators"
)

func (c *Client) GetEdgeGateway(ctx context.Context, params ParamsEdgeGateway) (*ModelEdgeGateway, error) {
	logger := c.logger.WithGroup("GetEdgeGateway")

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
	ep, _ := cav.GetEndpoint("EdgeGateway", cav.MethodGET)

	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.ID))
	if err != nil {
		logger.Error("Failed to get edge gateway", "error", err)
		return nil, err
	}

	return resp.Result().(*apiResponseEdgegateway).toModel(), nil
}

// TODO List,Create,Delete
