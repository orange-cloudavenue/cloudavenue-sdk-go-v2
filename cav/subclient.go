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

type subClientName string

const (
	ClientVmware    subClientName = "vmware"
	ClientCerberus  subClientName = "cerberus"
	ClientOSE       subClientName = "ose"
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

	ContextData(ctx context.Context) ContextData

	// getID returns stable cache identifier for subclient.
	getID() string

	close() error
}

var subClients = map[subClientName]subClientInterface{
	ClientVmware:    newVmwareClient(),
	ClientCerberus:  newCerberusClient(),
	ClientOSE:       newOSEClient(),
	ClientNetbackup: newNetbackupClient(),
}

func (s *subclient) getCredential() auth {
	return s.credential
}

func (s *subclient) setCredential(a auth) {
	s.credential = a
}

func (s *subclient) setConsole(console consoles.ConsoleName) {
	s.console = console
}

func (s *subclient) close() error {
	return nil
}
