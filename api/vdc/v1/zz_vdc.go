package vdc

import "context"

func (c *Client) ListVDC(ctx context.Context, params ParamsListVDC) (*ModelListVDC, error) {
	x, err := cmds.Get("VDC", "", "List").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelListVDC), nil
}

func (c *Client) GetVDC(ctx context.Context, params ParamsGetVDC) (*ModelGetVDC, error) {
	x, err := cmds.Get("VDC", "", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelGetVDC), nil
}

func (c *Client) CreateVDC(ctx context.Context, params ParamsCreateVDC) (*ModelGetVDC, error) {
	x, err := cmds.Get("VDC", "", "Create").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelGetVDC), nil
}

func (c *Client) UpdateVDC(ctx context.Context, params ParamsUpdateVDC) error {
	_, err := cmds.Get("VDC", "", "Update").Run(ctx, c, params)
	return err
}

func (c *Client) DeleteVDC(ctx context.Context, params ParamsDeleteVDC) error {
	_, err := cmds.Get("VDC", "", "Delete").Run(ctx, c, params)
	return err
}
