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
	"context"
)

// * Context

type contextKey string

const (
	contextKeyClientName contextKey = "subclient.clientName" // Context key for the client name
	contextExtraData     contextKey = "subclient.extraData"  // Context key for extra data
)

type ContextData struct {
	OrganizationID string // Organization ID for the context
	SiteID         string // Site ID for the context
}

// storeExtraDataInContext stores the extra data in the context.
func storeExtraDataInContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, contextExtraData, data)
}

// GetExtraDataFromContext retrieves the extra data from the context.
func GetExtraDataFromContext(ctx context.Context) ContextData {
	if data, ok := ctx.Value(contextExtraData).(ContextData); ok {
		return data
	}
	return ContextData{}
}
