package edgegateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/common-go/validators"
)

// ListT0 lists all T0s available in the organization.
func (c *Client) ListT0(ctx context.Context) (*ModelT0s, error) {
	// Get the endpoint for listing T0
	ep, _ := cav.GetEndpoint("T0", cav.MethodGET)

	// Perform the request to list T0
	resp, err := c.c.Do(ctx, ep)
	if err != nil {
		return nil, fmt.Errorf("error listing T0: %w", err)
	}

	return resp.Result().(*apiResponseT0s).ToModel(), nil
}

// GetTO retrieves a specific T0.
func (c *Client) GetT0(ctx context.Context, params ParamsGetT0) (*ModelT0, error) {
	if err := validators.New().Struct(&params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Get the endpoint for getting a specific T0
	ep, _ := cav.GetEndpoint("T0", cav.MethodGET)

	// Perform the request to get the specific T0
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], params.Name), // Only for mock response
	)
	if err != nil {
		return nil, fmt.Errorf("error getting T0: %w", err)
	}

	t0s := resp.Result().(*apiResponseT0s).ToModel()
	var t0 *ModelT0

	for _, t := range t0s.T0s {
		if t.Name == params.Name {
			t0 = &t
			break
		}
	}

	if t0 == nil {
		return nil, &errors.APIError{
			Operation:     "GetT0",
			StatusCode:    http.StatusNotFound,
			StatusMessage: http.StatusText(http.StatusNotFound),
			Message:       fmt.Sprintf("T0 with name %s not found", params.Name),
			Duration:      resp.Duration(),
			Endpoint:      resp.Request.URL,
			Method:        resp.Request.Method,
		}
	}

	return t0, nil
}

// GetT0FromEdgeGateway retrieves the T0 associated with a specific edge gateway.
func (c *Client) GetT0FromEdgeGateway(ctx context.Context, params ParamsEdgeGateway) (*ModelT0, error) {
	logger := c.logger.WithGroup("GetT0FromEdgeGateway")

	// Validate parameters
	if err := validators.New().Struct(&params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Get the endpoint for the edge gateway
	ep, _ := cav.GetEndpoint("T0", cav.MethodGET)

	// Perform the request to get the specific T0
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[1], params.Name), // Only for mock response
		cav.WithQueryParam(ep.QueryParams[2], params.ID),
	)
	if err != nil {
		logger.Error("Failed to get edge gateway", "error", err)
		return nil, err
	}

	t0s := resp.Result().(*apiResponseT0s).ToModel()
	for _, t0 := range t0s.T0s {
		for _, edgeGateway := range t0.EdgeGateways {
			if edgeGateway.ID == params.ID || edgeGateway.Name == params.Name {
				return &t0, nil
			}
		}
	}

	return nil, &errors.APIError{
		Operation:     "GetT0FromEdgeGateway",
		StatusCode:    http.StatusNotFound,
		StatusMessage: http.StatusText(http.StatusNotFound),
		Message: func() string {
			if params.ID != "" {
				return fmt.Sprintf("T0 for edge gateway with ID %s not found", params.ID)
			}
			return fmt.Sprintf("T0 for edge gateway with name %s not found", params.Name)
		}(),
		Duration: resp.Duration(),
		Method:   resp.Request.Method,
	}
}
