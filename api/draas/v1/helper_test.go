/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package draas

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
	assert.Nil(t, err, "Error creating draas client")
	return eC
}
