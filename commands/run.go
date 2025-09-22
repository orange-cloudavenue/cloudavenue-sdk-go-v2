/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

import (
	"context"
	"errors"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"
)

func (c *Command) Run(ctx context.Context, client, params any) (any, error) {
	c.params = params

	// If PreParamsRunnerFunc is defined, call it
	if c.PreParamsRunnerFunc != nil {
		paramsOut, err := c.PreParamsRunnerFunc(ctx, c, client, c.params)
		if err != nil {
			return nil, err
		}
		c.params = paramsOut
	}

	if len(c.ParamsSpecs) > 0 {
		if err := c.ParamsSpecs.buildAndValidateDynamicStruct(c.params); err != nil {
			return nil, err
		}
		// if err := c.ParamsSpecs.validate(c.params); err != nil {
		// 	return nil, err
		// }
	}

	// If PreRulesRunnerFunc is defined, call it
	if c.PreRulesRunnerFunc != nil {
		paramsOut, err := c.PreRulesRunnerFunc(ctx, c, client, c.params)
		if err != nil {
			return nil, err
		}
		c.params = paramsOut
	}

	if c.ParamsRules != nil {
		xlog.GetGlobalLogger().DebugContext(ctx, "Validating command params rules")
		cavClient, ok := getCavClientFromInterface(client)
		if !ok {
			return nil, errors.New("client must implement cav.Client interface")
		}
		if err := c.ParamsRules.validate(cavClient, c.params); err != nil {
			return nil, err
		}
	}

	v, err := c.RunnerFunc(ctx, c, client, params)
	if err != nil {
		return nil, err
	}
	return v, nil
}
