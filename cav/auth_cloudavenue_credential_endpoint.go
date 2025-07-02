package cav

import (
	"context"

	"resty.dev/v3"
)

func init() {
	err := Endpoint{
		Category:         CategoryAuthentication,
		Version:          VersionV1,
		Name:             "CreateSessionVmware",
		Method:           MethodPOST,
		SubClient:        ClientVmware,
		PathTemplate:     "/1.0.0/sessions",
		PathParams:       []PathParam{},
		QueryParams:      []QueryParam{},
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/sessions/post/",
		RequestFunc:      nil,
		requestInternalFunc: func(ctx context.Context, client *resty.Client, endpoint *Endpoint, opts ...RequestOption) (*resty.Response, error) {
			r := client.R().
				SetContext(ctx).
				SetHeader("Accept", "application/json;version="+VDCVersion)

			for _, opt := range opts {
				if err := opt(endpoint, r); err != nil {
					return nil, err
				}
			}

			return r.Post(endpoint.PathTemplate)
		},
	}.Register()
	if err != nil {
		panic(err)
	}
}
