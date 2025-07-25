package edgegateway

import "context"

func (c *Client) GetBandwidth(ctx context.Context, params ParamsEdgeGateway) (*ModelEdgeGatewayBandwidth, error) {
	x, err := cmds.Get("EdgeGateway", "Bandwidth", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelEdgeGatewayBandwidth), nil
}
