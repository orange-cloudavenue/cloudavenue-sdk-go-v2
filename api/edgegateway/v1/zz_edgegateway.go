package edgegateway

import "context"

func (c *Client) GetEdgeGateway(ctx context.Context, params ParamsEdgeGateway) (*ModelEdgeGateway, error) {
	x, err := cmds.Get("EdgeGateway", "", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelEdgeGateway), nil
}

func (c *Client) ListEdgeGateway(ctx context.Context) (*ModelEdgeGateways, error) {
	x, err := cmds.Get("EdgeGateway", "", "List").Run(ctx, c, nil)
	if err != nil {
		return nil, err
	}
	return x.(*ModelEdgeGateways), nil
}

func (c *Client) CreateEdgeGateway(ctx context.Context, params ParamsCreateEdgeGateway) (*ModelEdgeGateway, error) {
	x, err := cmds.Get("EdgeGateway", "", "Create").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelEdgeGateway), nil
}
