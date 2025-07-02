package org

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

type (
	Org struct {
		c cav.Client
	}

	Client interface{}
)

// New creates a new organization client.
func New(c cav.Client) (*Org, error) {
	if c == nil {
		return nil, errors.ErrClientNotInitialized
	}

	return &Org{
		c: c,
	}, nil
}
