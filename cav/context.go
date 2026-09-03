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
	"iter"
	"reflect"
)

type contextKey string

const (
	contextKeyClientName contextKey = "subclient.clientName"
	contextExtraData     contextKey = "subclient.extraData"
)

// ContextData carries backend-specific request context.
type ContextData struct {
	OrganizationID string // OrganizationID is CloudAvenue organization identifier.
	SiteID         string // SiteID is CloudAvenue site identifier.
}

// GetExtraDataFromContext returns backend-specific data stored in ctx.
func GetExtraDataFromContext(ctx context.Context) ContextData {
	if data, ok := ctx.Value(contextExtraData).(ContextData); ok {
		return data
	}
	return ContextData{}
}

// FieldsOf returns an iterator over P's exported struct fields.
// It exposes field metadata, including tags, for SDK consumers that need reflection.
func FieldsOf[P any]() iter.Seq[reflect.StructField] {
	return func(yield func(reflect.StructField) bool) {
		t := reflect.TypeFor[P]()
		for field := range t.Fields() {
			field := field
			if !field.IsExported() {
				continue
			}
			if !yield(field) {
				return
			}
		}
	}
}

// storeExtraDataInContext stores backend-specific data in ctx.
func storeExtraDataInContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, contextExtraData, data)
}
