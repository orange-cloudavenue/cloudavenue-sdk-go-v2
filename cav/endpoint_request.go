package cav

import (
	"context"

	"resty.dev/v3"
)

var DefaultRequestFunc = func(ctx context.Context, client Client, endpoint *Endpoint, opts ...RequestOption) (*resty.Response, error) {
	req, err := client.NewRequest(ctx, endpoint.SubClient)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(endpoint, req); err != nil {
			return nil, err
		}
	}
	return req.
		SetResult(endpoint.BodyType).
		Get(endpoint.PathTemplate)
}
