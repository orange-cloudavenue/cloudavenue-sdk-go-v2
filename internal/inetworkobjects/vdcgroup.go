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

type VDCGroupRef struct {
	ID   string
	Name string
}

type ResolveVDCGroupFunc func(ctx context.Context, id, name string) (VDCGroupRef, error)

func ResolveVDCGroupRef(ctx context.Context, id, name string, resolve ResolveVDCGroupFunc) (VDCGroupRef, error) {
	if id != "" && name != "" {
		return VDCGroupRef{ID: id, Name: name}, nil
	}

	return resolve(ctx, id, name)
}
