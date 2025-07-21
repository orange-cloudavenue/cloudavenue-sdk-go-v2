package edgegateway

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
)

func newClient(t *testing.T) *Client {
	t.Helper()

	mC, err := mock.NewClient(
		mock.WithLogger(
			slog.New(
				slog.NewTextHandler(
					os.Stdout,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					}),
			),
		),
	)
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")
	return eC
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
