package edgegateway

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

func (c *Client) retrieveEdgeGatewayIDByName(ctx context.Context, name string) (string, error) {
	epQuery, _ := cav.GetEndpoint("QueryEdgeGateway", cav.MethodGET)

	respQuery, err := c.c.Do(
		ctx,
		epQuery,
		cav.WithQueryParam(epQuery.QueryParams[1], "name=="+name),
	)
	if err != nil {
		return "", err
	}

	// Record is already checked in the middleware.
	return respQuery.Result().(*apiResponseQueryEdgeGateway).Record[0].ID, nil
}
