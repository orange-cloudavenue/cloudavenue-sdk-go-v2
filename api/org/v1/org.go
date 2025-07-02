package org

import (
	"context"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

func init() {
	err := cav.Endpoint{
		Category:     "demo",
		Version:      cav.VersionV1,
		Name:         "demo-api",
		Method:       cav.MethodGET,
		SubClient:    cav.ClientVmware,
		PathTemplate: "/1.0.0/orgs/{orgUrn}",
		PathParams: []cav.PathParam{
			{
				Name:        "orgUrn",
				Description: "The organization URN (Uniform Resource Name) to identify the organization.",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn_rfc2141")
				},
			},
		},
		QueryParams:      []cav.QueryParam{},
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/orgs/orgUrn/get/",
		BodyType:         OrgResponse{},
	}.Register()

	if err != nil {
		panic(err)
	}
}

type OrgResponse struct { //nolint:revive
	CanManageOrgs bool   `json:"canManageOrgs" `
	CanPublish    bool   `json:"canPublish"`
	CatalogCount  int64  `json:"catalogCount"`
	Description   string `json:"description"`
	DiskCount     int64  `json:"diskCount"`
	DisplayName   string `json:"displayName"`
	ID            string `json:"id"`
	IsEnabled     bool   `json:"isEnabled"`
	ManagedBy     struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"managedBy"`
	Name           string `json:"name"`
	OrgVdcCount    int64  `json:"orgVdcCount"`
	RunningVMCount int64  `json:"runningVMCount"`
	UserCount      int64  `json:"userCount"`
	VappCount      int64  `json:"vappCount"`
}

// DemoRequest represents a request to the demo cav.
func DemoRequest(ctx context.Context, client cav.Client, orgID string) (*OrgResponse, error) {
	demoEndpoint, err := cav.GetEndpoint("demo", cav.VersionV1, "demo-api", cav.MethodGET)
	if err != nil {
		return nil, err
	}

	resp, err := demoEndpoint.RequestFunc(
		ctx,
		client,
		demoEndpoint,
		cav.WithPathParam(demoEndpoint.PathParams[0], orgID),
	)
	if err != nil {
		return nil, err
	}

	if err := client.ParseAPIError("Get organization detail", resp); err != nil {
		return nil, err
	}

	return resp.Result().(*OrgResponse), nil
}
