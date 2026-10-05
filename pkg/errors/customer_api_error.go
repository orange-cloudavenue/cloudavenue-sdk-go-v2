/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package errors

import (
	"net/http"
	"regexp"
	"strings"
	"time"
)

// CustomerAPIErrorBody is error body returned by customer API.
//
// Live body shape:
//
//	{"minorErrorCode":"BAD_REQUEST","message":"...","stackTrace":null}
//
// Code and Reason are fallback fields for older or alternate payloads.
type CustomerAPIErrorBody struct {
	MinorErrorCode string `json:"minorErrorCode" fake:"{regex:[A-Z_]{8,20}}"`
	Code           string `json:"code" fake:"{regex:err-[0-9]{4}}"`
	Reason         string `json:"reason" fake:"{regex:mock-[0-9]{4}}"`
	Message        string `json:"message" fake:"{sentence:3,10}"`
}

const customerAPIWAFRejectedCode = "REQUEST_REJECTED"

var regexCustomerAPIWAFRejected = regexp.MustCompile(`(?i)<title>\s*Request Rejected\s*</title>|The requested URL was rejected`)

var regexCustomerAPIWAFSupportID = regexp.MustCompile(`support ID is:\s*(\d+)`)

// ErrorCode returns symbolic code from body, preferring live field.
func (e *CustomerAPIErrorBody) ErrorCode() string {
	if e == nil {
		return ""
	}
	if e.MinorErrorCode != "" {
		return e.MinorErrorCode
	}
	return e.Code
}

// ErrorMessage returns human-readable message built from body fields.
func (e *CustomerAPIErrorBody) ErrorMessage(unknownMessage string) string {
	if e == nil {
		return unknownMessage
	}

	prefix := e.ErrorCode()
	if prefix == "" {
		prefix = e.Reason
	}

	switch {
	case prefix == "" && e.Message == "":
		return unknownMessage
	case prefix == "":
		return e.Message
	case e.Message == "":
		return prefix
	default:
		return prefix + ": " + e.Message
	}
}

// CustomerAPIStatusError builds APIError from status, payload, transport metadata.
func CustomerAPIStatusError(operation string, statusCode int, body *CustomerAPIErrorBody, unknownMessage string, duration time.Duration, endpoint, method string, sentinel error) *APIError {
	if body != nil {
		return &APIError{
			Operation:     operation,
			StatusCode:    statusCode,
			StatusMessage: body.ErrorCode(),
			Message:       body.ErrorMessage(unknownMessage),
			Duration:      duration,
			Endpoint:      endpoint,
			Method:        method,
			Err:           sentinel,
		}
	}

	return &APIError{
		Operation:  operation,
		StatusCode: statusCode,
		Message:    unknownMessage,
		Duration:   duration,
		Endpoint:   endpoint,
		Method:     method,
		Err:        sentinel,
	}
}

// CustomerAPIWAFError returns APIError for WAF rejection page served as success.
func CustomerAPIWAFError(operation string, statusCode int, contentType, rawBody string, duration time.Duration, endpoint, method string) *APIError {
	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		return nil
	}

	if !regexCustomerAPIWAFRejected.MatchString(rawBody) {
		return nil
	}

	msg := "request rejected by web application firewall"
	if m := regexCustomerAPIWAFSupportID.FindStringSubmatch(rawBody); len(m) == 2 {
		msg += " (support ID: " + m[1] + ")"
	}

	return &APIError{
		Operation:     operation,
		StatusCode:    statusCode,
		StatusMessage: customerAPIWAFRejectedCode,
		Message:       msg,
		Duration:      duration,
		Endpoint:      endpoint,
		Method:        method,
		Err:           classifyHTTPStatusCode(statusCode),
	}
}

// CustomerAPIWAFRejectedCode returns machine-readable code for WAF rejections.
func CustomerAPIWAFRejectedCode() string {
	return customerAPIWAFRejectedCode
}

func classifyHTTPStatusCode(statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusMethodNotAllowed:
		return ErrMethodNotAllowed
	case http.StatusRequestTimeout:
		return ErrRequestTimeout
	case http.StatusConflict:
		return ErrConflict
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusInternalServerError:
		return ErrInternalServerError
	case http.StatusBadGateway:
		return ErrBadGateway
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	case http.StatusGatewayTimeout:
		return ErrGatewayTimeout
	default:
		return nil
	}
}
