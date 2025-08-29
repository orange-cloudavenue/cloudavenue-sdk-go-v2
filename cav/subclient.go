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

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var subClients = map[subClientName]subClientInterface{
	ClientVmware:   newVmwareClient(),
	ClientCerberus: newCerberusClient(),
}

type subClientName string

const (
	ClientVmware    subClientName = "vmware"
	ClientCerberus  subClientName = "cerberus"
	ClientNetbackup subClientName = "netbackup"
)

type subclient struct {
	credential auth
	console    consoles.ConsoleName
}

type subClientInterface interface {
	setCredential(auth)
	getCredential() auth
	setConsole(consoles.ConsoleName)
	newHTTPClient(context.Context) (*resty.Client, error)

	parseAPIError(operation string, resp *resty.Response) *errors.APIError
	idempotentRetryCondition() resty.RetryConditionFunc

	// getID returns the unique identifier for the subclient
	getID() string

	close() error
}

// getCredential retrieves the current authentication credentials.
func (s *subclient) getCredential() auth {
	return s.credential
}

// setCredential sets the authentication credential for the subclient.
func (s *subclient) setCredential(a auth) {
	s.credential = a
}

// setConsole sets the console name for the subclient.
func (s *subclient) setConsole(console consoles.ConsoleName) {
	s.console = console
}
