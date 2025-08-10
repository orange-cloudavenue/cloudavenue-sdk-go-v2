package commands

import (
	"context"
	"errors"
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
		if err := c.ParamsSpecs.validate(params); err != nil {
			return nil, err
		}
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
