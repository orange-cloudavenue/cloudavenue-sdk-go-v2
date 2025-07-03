/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

type (
	Org struct {
		c cav.Client
	}
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
