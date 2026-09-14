/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"

func classifyStatusCode(statusCode int) error {
	switch statusCode {
	case 400:
		return errors.ErrBadRequest
	case 401:
		return errors.ErrUnauthorized
	case 403:
		return errors.ErrForbidden
	case 404:
		return errors.ErrNotFound
	case 405:
		return errors.ErrMethodNotAllowed
	case 408:
		return errors.ErrRequestTimeout
	case 409:
		return errors.ErrConflict
	case 429:
		return errors.ErrTooManyRequests
	case 500:
		return errors.ErrInternalServerError
	case 502:
		return errors.ErrBadGateway
	case 503:
		return errors.ErrServiceUnavailable
	case 504:
		return errors.ErrGatewayTimeout
	default:
		return nil
	}
}
