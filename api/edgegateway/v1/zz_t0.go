package edgegateway

import "context"

func (c *Client) ListT0(ctx context.Context) (*ModelT0s, error) {
	x, err := cmds.Get("T0", "", "List").Run(ctx, c, nil)
	if err != nil {
		return nil, err
	}
	return x.(*ModelT0s), nil
}

func (c *Client) GetT0(ctx context.Context, params ParamsGetT0) (*ModelT0, error) {
	x, err := cmds.Get("T0", "", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelT0), nil
}
