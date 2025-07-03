package cav

import (
	"testing"
)

type mockAPI struct{}
type mockVersion struct{}

func TestWithExtraProperties(t *testing.T) {
	var (
		api  = API("mock")
		opts endpointRegistryOptions
	)
	optFunc := WithExtraProperties(api, VersionV1)
	optFunc(&opts)
	if opts.api != api {
		t.Errorf("expected api to be set")
	}
	if opts.version != VersionV1 {
		t.Errorf("expected version to be set")
	}
}
