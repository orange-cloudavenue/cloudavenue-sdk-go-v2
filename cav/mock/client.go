/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package mock

import (
	"log"
	"log/slog"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

const (
	mockOrg = "cav01ev01ocb0001234"
)

var logger = xlog.GetGlobalLogger()

func NewClient(opts ...OptionFunc) (cav.Client, *MockServer, error) {
	Options := &Options{}
	for _, opt := range opts {
		if err := opt(Options); err != nil {
			return nil, nil, err
		}
	}

	if Options.logger != nil {
		xlog.SetGlobalLogger(Options.logger)
		logger = Options.logger
	}

	ms := newServer(logger)

	// Get all endpoints and register handlers on the chi router.
	// Handlers are registered on the real path templates, so no URL rewriting is needed.
	endpoints := cav.GetEndpointsUncategorized()
	mux := chi.NewRouter()

	// Group endpoints by (method, path) to handle path collisions.
	groups := make(map[string][]*cav.Endpoint)
	for _, ep := range endpoints {
		key := ep.Method.String() + "|" + ep.PathTemplate
		groups[key] = append(groups[key], ep)
	}

	for _, group := range groups {
		ep := group[0]
		logger.Debug(
			"Registering mock endpoint",
			slog.String("name", ep.Name),
			slog.String("method", ep.Method.String()),
			slog.String("path", ep.PathTemplate),
			slog.String("ID", ep.ID),
			slog.Int("count", len(group)),
		)
		if len(group) == 1 {
			mux.MethodFunc(ep.Method.String(), ep.PathTemplate, ms.handlerFor(ep))
		} else {
			mux.MethodFunc(ep.Method.String(), ep.PathTemplate, ms.handlerForGroup(group))
		}
	}

	// Pre-register the SessionVmware mock handler.
	if ep, err := cav.GetEndpoint("SessionVmware"); err == nil {
		ms.SetResponseFunc(ep, cav.NewSessionVmwareMockHandler())
	}

	hts := httptest.NewServer(mux)
	slog.SetDefault(logger)
	hts.Config.ErrorLog = log.Default()

	logger.Debug("Mock server created", slog.String("url", hts.URL))

	nC, err := cav.NewClient(
		mockOrg,
		cav.WithCustomEndpoints(consoles.Services{
			IHM: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/ihm",
			},
			APIVCD: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL,
			},
			APICerberus: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL,
			},
			S3: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/s3",
			},
			Netbackup: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/netbackup",
			},
		}),
		cav.WithCloudAvenueCredential("mockuser", "mockpassword"),
		cav.WithLogger(logger),
	)
	if err != nil {
		return nil, nil, err
	}

	logger.Debug("Mock client created", slog.String("organization", mockOrg))

	return nC, ms, nil
}

var GetEndpoint = cav.GetEndpoint
