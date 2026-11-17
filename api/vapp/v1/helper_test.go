/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
)

func newClient(t *testing.T) (*Client, *mock.Server) {
	t.Helper()

	mC, ms, err := mock.NewClient(
		mock.WithLogger(
			slog.New(
				slog.NewTextHandler(
					os.Stdout,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			),
		),
	)
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating vapp client")
	return eC, ms
}

func xmlResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("X-Cloud-Avenue-Mock", "true")
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TestClientJobOptions(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		mC, _, err := mock.NewClient()
		assert.NoError(t, err)

		client, err := New(mC)
		assert.NoError(t, err)
		if assert.NotNil(t, client) {
			assert.Equal(t, defaultJobTimeout, client.jobTimeout)
			assert.Equal(t, defaultJobPollInterval, client.jobPollInterval)
		}
	})

	t.Run("options applied", func(t *testing.T) {
		mC, _, err := mock.NewClient()
		assert.NoError(t, err)

		client, err := New(mC, WithJobTimeout(2*time.Minute), WithJobPollInterval(500*time.Millisecond))
		assert.NoError(t, err)
		if assert.NotNil(t, client) {
			assert.Equal(t, 2*time.Minute, client.jobTimeout)
			assert.Equal(t, 500*time.Millisecond, client.jobPollInterval)
		}
	})

	t.Run("rejects non-positive values", func(t *testing.T) {
		mC, _, err := mock.NewClient()
		assert.NoError(t, err)

		client, err := New(mC, WithJobTimeout(0))
		assert.Error(t, err)
		assert.Nil(t, client)

		client, err = New(mC, WithJobTimeout(-time.Second))
		assert.Error(t, err)
		assert.Nil(t, client)

		client, err = New(mC, WithJobPollInterval(0))
		assert.Error(t, err)
		assert.Nil(t, client)

		client, err = New(mC, WithJobPollInterval(-time.Second))
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}
