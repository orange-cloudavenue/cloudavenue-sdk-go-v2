/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opUndeployVApp              = "VApp.Undeploy"
	opDeployVApp                = "VApp.Deploy"
	opPowerOffVApp              = "VApp.PowerOff"
	opPowerOnVApp               = "VApp.PowerOn"
	opRebootVApp                = "VApp.Reboot"
	opResetVApp                 = "VApp.Reset"
	opShutdownVApp              = "VApp.Shutdown"
	opSuspendVApp               = "VApp.Suspend"
	opDiscardSuspendedStateVApp = "VApp.DiscardSuspendedState"
)

var (
	undeployVAppOp = cav.Operation[undeployVAppParams, cav.Job]{
		Name:     opUndeployVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UndeployVApp(),
		Validate: func(p undeployVAppParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p undeployVAppParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UndeployVApp()
			return []cav.EndpointRequestOption{
				cav.WithPathParam(ep.PathParams[0], p.ID),
				cav.SetCustomRestyOption(func(req *resty.Request) {
					req.SetHeader("Content-Type", "application/xml")
				}),
			}, nil
		},
		Transform: func(p undeployVAppParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ undeployVAppParams) (cav.Job, error) {
			return extractVMwareJob(resp, "undeploy")
		},
	}
	deployVAppOp                = newVAppActionOperation(opDeployVApp, endpoints.DeployVApp, "deploy")
	powerOffVAppOp              = newVAppActionOperation(opPowerOffVApp, endpoints.PowerOffVApp, "power off")
	powerOnVAppOp               = newVAppActionOperation(opPowerOnVApp, endpoints.PowerOnVApp, "power on")
	rebootVAppOp                = newVAppActionOperation(opRebootVApp, endpoints.RebootVApp, "reboot")
	resetVAppOp                 = newVAppActionOperation(opResetVApp, endpoints.ResetVApp, "reset")
	shutdownVAppOp              = newVAppActionOperation(opShutdownVApp, endpoints.ShutdownVApp, "shutdown")
	suspendVAppOp               = newVAppActionOperation(opSuspendVApp, endpoints.SuspendVApp, "suspend")
	discardSuspendedStateVAppOp = newVAppActionOperation(opDiscardSuspendedStateVApp, endpoints.DiscardSuspendedStateVApp, "discard suspended state")
)

func newVAppActionOperation(name string, endpoint func() *cav.Endpoint, action string) cav.Operation[getVAppByIDParams, cav.Job] {
	return cav.Operation[getVAppByIDParams, cav.Job]{
		Name:     name,
		Backend:  cav.BackendVMware,
		Endpoint: endpoint(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoint()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (cav.Job, error) {
			return extractVMwareJob(resp, action)
		},
	}
}

// DeployVApp deploys a VApp by ID.
func (c *Client) DeployVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, deployVAppOp, opDeployVApp)
}

// UndeployVApp undeploys a VApp using the requested VMware undeploy mode.
func (c *Client) UndeployVApp(ctx context.Context, params types.ParamsUndeployVApp) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opUndeployVApp, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return fmt.Errorf("%s: %w", opUndeployVApp, err)
	}

	job, err := cav.Execute(ctx, c.c, undeployVAppOp, undeployVAppParams{ID: vappID, Body: itypes.ToAPIRequestUndeployVApp(params.Mode)})
	if err != nil {
		return fmt.Errorf("%s: %w", opUndeployVApp, err)
	}
	if err := awaitVAppJob(ctx, c.c, opUndeployVApp, job.ID); err != nil {
		return err
	}
	return nil
}

// PowerOffVApp powers off a VApp by ID.
func (c *Client) PowerOffVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, powerOffVAppOp, opPowerOffVApp)
}

// PowerOnVApp powers on a VApp by ID.
func (c *Client) PowerOnVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, powerOnVAppOp, opPowerOnVApp)
}

// RebootVApp reboots a VApp by ID.
func (c *Client) RebootVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, rebootVAppOp, opRebootVApp)
}

// ResetVApp resets a VApp by ID.
func (c *Client) ResetVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, resetVAppOp, opResetVApp)
}

// ShutdownVApp shuts down a VApp by ID.
func (c *Client) ShutdownVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, shutdownVAppOp, opShutdownVApp)
}

// SuspendVApp suspends a VApp by ID.
func (c *Client) SuspendVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, suspendVAppOp, opSuspendVApp)
}

// DiscardSuspendedStateVApp discards the suspended state of a VApp by ID.
func (c *Client) DiscardSuspendedStateVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, discardSuspendedStateVAppOp, opDiscardSuspendedStateVApp)
}

func (c *Client) executeVAppAction(ctx context.Context, id string, op cav.Operation[getVAppByIDParams, cav.Job], name string) error {
	job, err := cav.Execute(ctx, c.c, op, getVAppByIDParams{ID: id})
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if err := awaitVAppJob(ctx, c.c, name, job.ID); err != nil {
		return err
	}
	return nil
}

func awaitVAppJob(ctx context.Context, c cav.Client, opName, jobID string) error {
	_, err := cav.AwaitJobOnBackend(ctx, c, cav.BackendVMware, jobID, cav.JobPollOptions{
		Timeout:         30 * time.Second,
		PollingInterval: 1 * time.Second,
	}, func(_ *resty.Response) (struct{}, error) {
		return struct{}{}, nil
	})
	if err != nil {
		return fmt.Errorf("%s: await job: %w", opName, err)
	}
	return nil
}

func isIgnorableUndeployError(err error) bool {
	var apiErr *pkgerrors.APIError
	if !stderrors.As(err, &apiErr) || apiErr == nil {
		return false
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		return false
	}

	message := strings.ToLower(strings.TrimSpace(apiErr.Message))
	return strings.Contains(message, "requested operation could not be executed") &&
		strings.Contains(message, "vapp") &&
		strings.Contains(message, "not running")
}
