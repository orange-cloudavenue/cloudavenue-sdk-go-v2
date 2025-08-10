package vdc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewClient_ClientNil(t *testing.T) {
	c, err := New(nil)
	assert.Nil(t, c, "Expected nil client when input is nil")
	assert.Error(t, err, "Expected error when input is nil")
}
