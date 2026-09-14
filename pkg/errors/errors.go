/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package errors

var (
	ErrClientNotInitialized = New("client not initialized")
	ErrNotFound             = New("not found")
	ErrBadRequest           = New("bad request")
	ErrUnauthorized         = New("unauthorized")
	ErrForbidden            = New("forbidden")
	ErrMethodNotAllowed     = New("method not allowed")
	ErrRequestTimeout       = New("request timeout")
	ErrConflict             = New("conflict")
	ErrTooManyRequests      = New("too many requests")
	ErrInternalServerError  = New("internal server error")
	ErrBadGateway           = New("bad gateway")
	ErrServiceUnavailable   = New("service unavailable")
	ErrGatewayTimeout       = New("gateway timeout")
	ErrJobFailed            = New("job failed")
	ErrJobTimeout           = New("job timeout")
)
