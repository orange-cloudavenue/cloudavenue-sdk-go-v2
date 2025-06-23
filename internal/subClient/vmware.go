package subclient

import (
	"context"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ Client = &vmware{}

type vmware struct {
	client
}

type vmwareError struct {
	Message        string `json:"message"`
	MinorErrorCode string `json:"minorErrorCode"`
}

var NewVmwareClient = func() Client {
	return &vmware{}
}

// NewClient creates a new request for the VMware subclient.
func (v *vmware) NewHTTPClient(ctx context.Context) (*resty.Client, error) {
	if !v.credential.IsInitialized() {
		if err := v.credential.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	v.httpClient = httpclient.NewHTTPClient().
		SetBaseURL(v.console.GetAPIVmwareEndpoint()).
		SetHeaders(v.credential.Headers()).
		SetHeader("Accept", "application/json;version="+auth.VDCVersion).
		SetError(vmwareError{})

	return v.httpClient, nil
}

// SetCredential sets the authentication credential for the VMware client.
func (v *vmware) SetCredential(a auth.Auth) {
	v.credential = a
}

// SetConsole sets the console for the VMware client.
func (v *vmware) SetConsole(c consoles.Console) {
	v.console = c
}

// ParseAPIError parses the API error response from the VMware client.
func (v *vmware) ParseAPIError(resp *resty.Response) *errors.APIError {
	if resp == nil || !resp.IsError() {
		return nil
	}

	// If resp.Error() is not nil, it means an error occurred.
	// Parse the error response body.
	if err, ok := resp.Error().(*vmwareError); ok {
		return &errors.APIError{
			StatusCode: resp.StatusCode(),
			Message:    err.Message,
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
		}
	}

	return nil
}
