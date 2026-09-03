/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package inetworkobjects

import "context"

type VdcGroupRef struct {
	ID   string
	Name string
}

type ResolveVdcGroupFunc func(ctx context.Context, id, name string) (VdcGroupRef, error)

func ResolveVdcGroupRef(ctx context.Context, id, name string, resolve ResolveVdcGroupFunc) (VdcGroupRef, error) {
	if id != "" && name != "" {
		return VdcGroupRef{ID: id, Name: name}, nil
	}

	return resolve(ctx, id, name)
}
