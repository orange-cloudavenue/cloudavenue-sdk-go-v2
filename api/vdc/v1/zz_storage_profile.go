package vdc

import "context"

func (c *Client) ListStorageProfile(ctx context.Context, params ParamsListStorageProfiles) (*ModelListStorageProfiles, error) {
	x, err := cmds.Get("VDC", "StorageProfile", "List").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelListStorageProfiles), nil
}

func (c *Client) AddStorageProfile(ctx context.Context, params ParamsAddStorageProfile) error {
	_, err := cmds.Get("VDC", "StorageProfile", "Add").Run(ctx, c, params)
	return err
}
