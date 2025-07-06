/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"

// logger can be overridden by WithLogger option.
// If not set, it defaults to xlog.New.
// It is used for logging messages in the client.
var xlogger = xlog.GetGlobalLogger()
