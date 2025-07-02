package org

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

func TestOrgNew_WithNilClient(t *testing.T) {
	orgClient, err := New(nil)
	assert.Nil(t, orgClient)
	assert.Equal(t, errors.ErrClientNotInitialized, err)
}

func TestOrgNew_WithValidClient(t *testing.T) {
	mockClient, err := mock.NewClient()
	assert.Nil(t, err)

	orgClient, err := New(mockClient)
	assert.NotNil(t, orgClient)
	assert.Nil(t, err)
}
