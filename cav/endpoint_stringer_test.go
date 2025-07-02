package cav

import (
	"testing"
)

func TestEndpoint_String(t *testing.T) {
	e := Endpoint{
		Category:     "cat",
		Version:      "v1",
		Name:         "endpointName",
		Method:       "GET",
		PathTemplate: "/path/{id}",
	}
	expected := "[cat] v1 endpointName GET /path/{id}"
	if got := e.String(); got != expected {
		t.Errorf("Endpoint.String() = %q, want %q", got, expected)
	}
}
