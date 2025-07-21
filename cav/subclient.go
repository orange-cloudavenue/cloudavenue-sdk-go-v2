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

var subClients = map[SubClientName]SubClient{
	ClientVmware:   newVmwareClient(),
	ClientCerberus: newCerberusClient(),
}

type SubClientName string

const (
	ClientVmware    SubClientName = "vmware"
	ClientCerberus  SubClientName = "cerberus"
	ClientNetbackup SubClientName = "netbackup"
)

type subclient struct {
	httpClient *resty.Client
	credential auth
	console    consoles.Console
}

type SubClient interface {
	SetCredential(auth)
	SetConsole(consoles.Console)
	NewHTTPClient(context.Context) (*resty.Client, error)
	parseAPIError(operation string, resp *resty.Response) *errors.APIError
	idempotentRetryCondition() resty.RetryConditionFunc
}
