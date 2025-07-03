/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
