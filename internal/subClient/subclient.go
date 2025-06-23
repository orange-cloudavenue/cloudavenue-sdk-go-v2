package subclient

import (
	"context"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var Clients = map[Name]Client{
	Vmware:   NewVmwareClient(),
	Cerberus: NewCerberusClient(),
	// Mock client for testing purposes
	mock: NewMockClient(),
}

type Name string

const (
	Vmware    Name = "vmware"
	Cerberus  Name = "cerberus"
	Netbackup Name = "netbackup"
	mock      Name = "mock" // For testing purposes
)

type client struct {
	httpClient *resty.Client
	credential auth.Auth
	console    consoles.Console
}

type Client interface {
	SetCredential(auth.Auth)
	SetConsole(consoles.Console)
	NewHTTPClient(context.Context) (*resty.Client, error)
	ParseAPIError(*resty.Response) *errors.APIError
}
