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
	"errors"
	"log/slog"
	"time"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

const (
	defaultJobTimeout      = 5 * time.Minute
	defaultJobPollInterval = 1 * time.Second
)

type (
	Client struct {
		c               cav.Client
		logger          *slog.Logger
		jobTimeout      time.Duration
		jobPollInterval time.Duration
	}

	// Option configures the vapp client at construction time.
	Option func(*Client)
)

// WithJobTimeout sets the timeout applied when awaiting VMware vApp jobs.
func WithJobTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.jobTimeout = timeout
	}
}

// WithJobPollInterval sets the polling interval applied when awaiting VMware vApp jobs.
func WithJobPollInterval(interval time.Duration) Option {
	return func(c *Client) {
		c.jobPollInterval = interval
	}
}

// New creates a new vapp client.
func New(c cav.Client, opts ...Option) (*Client, error) {
	if c == nil {
		return nil, pkgerrors.ErrClientNotInitialized
	}

	client := &Client{
		c:               c,
		logger:          c.Logger().WithGroup("vapp"),
		jobTimeout:      defaultJobTimeout,
		jobPollInterval: defaultJobPollInterval,
	}
	for _, opt := range opts {
		opt(client)
	}

	if client.jobTimeout <= 0 {
		return nil, errors.New("vapp: job timeout must be positive")
	}
	if client.jobPollInterval <= 0 {
		return nil, errors.New("vapp: job poll interval must be positive")
	}

	client.logger.Debug("Successfully creating new client")

	return client, nil
}
