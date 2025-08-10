package vdc

import (
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
)

var testMutex = sync.Mutex{}

func newClient(t *testing.T) *Client {
	t.Helper()

	testMutex.Lock()
	t.Cleanup(func() {
		testMutex.Unlock()
	})

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
	assert.Nil(t, err, "Error creating vdc client")
	return eC
}
