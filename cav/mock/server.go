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
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"
	"sort"
	"sync"

	"github.com/orange-cloudavenue/common-go/generator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// mockResponse holds per-endpoint mock response configuration.
type mockResponse struct {
	data       any
	statusCode *int
	handler    http.HandlerFunc
}

// Server holds mock response state for all endpoints.
type Server struct {
	mu        sync.RWMutex
	responses map[string]*mockResponse
	logger    *slog.Logger
}

// newServer creates Server.
func newServer(logger *slog.Logger) *Server {
	return &Server{
		responses: make(map[string]*mockResponse),
		logger:    logger,
	}
}

// SetResponse sets mock response payload and status code for endpoint.
func (ms *Server) SetResponse(ep *cav.Endpoint, data any, statusCode *int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	mr, ok := ms.responses[ep.Name]
	if !ok {
		mr = &mockResponse{}
		ms.responses[ep.Name] = mr
	}

	mr.data = data
	mr.statusCode = statusCode
	mr.handler = nil

	ms.logger.Debug("Mock response set", slog.String("endpoint", ep.Name), slog.Any("status_code", statusCode))
}

// SetResponseFunc sets custom handler for endpoint.
func (ms *Server) SetResponseFunc(ep *cav.Endpoint, handler http.HandlerFunc) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	mr, ok := ms.responses[ep.Name]
	if !ok {
		mr = &mockResponse{}
		ms.responses[ep.Name] = mr
	}

	mr.handler = handler
	mr.data = nil
	mr.statusCode = nil

	ms.logger.Debug("Mock response func set", slog.String("endpoint", ep.Name))
}

// CleanResponse removes mock override for endpoint.
func (ms *Server) CleanResponse(ep *cav.Endpoint) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.responses, ep.Name)
	ms.logger.Debug("Mock response cleaned", slog.String("endpoint", ep.Name))
}

// handlerFor builds handler for endpoint.
func (ms *Server) handlerFor(ep *cav.Endpoint) http.HandlerFunc {
	if ep.ResponseType != nil {
		bodyType := reflect.TypeOf(ep.ResponseType)
		if bodyType.Kind() == reflect.Pointer {
			bodyType = bodyType.Elem()
		}
		if bodyType == reflect.TypeFor[cav.Job]() || bodyType == reflect.TypeFor[cav.CerberusJobCreatedAPIResponse]() {
			return ms.jobHandlerFor(ep)
		}
	}

	return ms.makeHandler(ep)
}

// makeHandler builds default request handler for endpoint.
func (ms *Server) makeHandler(ep *cav.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ms.mu.RLock()
		mr, ok := ms.responses[ep.Name]
		ms.mu.RUnlock()

		if ok && mr.handler != nil {
			ms.logger.Debug("Using custom mock handler", slog.String("endpoint", ep.Name))
			mr.handler(w, r)
			return
		}

		statusCode := http.StatusOK
		var data any

		if ok {
			if mr.statusCode != nil {
				statusCode = *mr.statusCode
			}
			data = mr.data
		}

		if statusCode >= 300 {
			ms.writeErrorResponse(w, ep, statusCode)
			return
		}

		if data != nil {
			ms.writeJSONResponse(w, ep, statusCode, data)
			return
		}

		if ep.ResponseType != nil {
			body := ms.generateBody(ep)
			if body == nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			ms.writeJSONResponse(w, ep, statusCode, body)
			return
		}

		ms.writeJSONResponse(w, ep, statusCode, map[string]string{"message": "Mock response"})
	}
}

// handlerForGroup routes colliding method/path pairs by query parameters.
func (ms *Server) handlerForGroup(eps []*cav.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ep := ms.matchEndpoint(r, eps)
		handler := ms.handlerFor(ep)
		handler(w, r)
	}
}

// matchEndpoint returns best matching endpoint for request.
func (ms *Server) matchEndpoint(r *http.Request, eps []*cav.Endpoint) *cav.Endpoint {
	sort.Slice(eps, func(i, j int) bool {
		return eps[i].Name < eps[j].Name
	})

	query := r.URL.Query()
	bestEp := eps[0]
	bestScore := -1

	for _, ep := range eps {
		score := 0
		for _, qp := range ep.QueryParams {
			if _, ok := query[qp.Name]; ok {
				score++
				if qp.Value != "" && query.Get(qp.Name) == qp.Value {
					score++
				}
			}
		}
		if score > bestScore || (score == bestScore && len(ep.QueryParams) < len(bestEp.QueryParams)) {
			bestScore = score
			bestEp = ep
		}
	}

	if bestScore == 0 {
		for _, ep := range eps {
			if ms.hasMock(ep) {
				return ep
			}
		}
	}

	return bestEp
}

// hasMock reports whether endpoint has explicit mock override.
func (ms *Server) hasMock(ep *cav.Endpoint) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	mr, ok := ms.responses[ep.Name]
	return ok && mr != nil && (mr.handler != nil || mr.data != nil || mr.statusCode != nil)
}

// writeErrorResponse writes error response for status codes >= 300.
func (ms *Server) writeErrorResponse(w http.ResponseWriter, ep *cav.Endpoint, statusCode int) {
	ms.logger.Debug("Mock error response", slog.String("endpoint", ep.Name), slog.Int("statusCode", statusCode))
	http.Error(w, http.StatusText(statusCode), statusCode)
}

// generateBody auto-generates response body from endpoint ResponseType.
func (ms *Server) generateBody(ep *cav.Endpoint) any {
	bodyType := reflect.TypeOf(ep.ResponseType)
	if bodyType.Kind() == reflect.Pointer {
		bodyType = bodyType.Elem()
	}

	if bodyType == reflect.TypeFor[cav.Job]() {
		return nil
	}

	if bodyType == reflect.TypeFor[cav.CerberusJobCreatedAPIResponse]() {
		return &cav.CerberusJobCreatedAPIResponse{
			ID:      "87ab1934-0146-4fb0-80bc-815fea03214d",
			Message: "Job created successfully",
		}
	}

	newBodyType := reflect.PointerTo(bodyType)
	newBody := reflect.New(newBodyType).Interface()

	switch bodyType.Kind() {
	case reflect.Slice:
		generator.Slice(newBody)
	default:
		if err := generator.Struct(newBody); err != nil {
			ms.logger.Error("Error generating mock data", slog.String("endpoint", ep.Name), slog.Any("error", err))
			return nil
		}
	}

	return newBody
}

// writeJSONResponse writes JSON response with status code.
func (ms *Server) writeJSONResponse(w http.ResponseWriter, ep *cav.Endpoint, statusCode int, data any) {
	bodyEncoded, err := json.Marshal(data)
	if err != nil {
		ms.logger.Error("Error encoding mock response", slog.String("endpoint", ep.Name), slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cloud-Avenue-Mock", "true")

	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	_, err = w.Write(bodyEncoded)
	if err != nil {
		ms.logger.Error("Error writing mock response", slog.String("endpoint", ep.Name), slog.Any("error", err))
	}
}

// jobHandlerFor builds handler for job-oriented endpoints.
func (ms *Server) jobHandlerFor(ep *cav.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ms.mu.RLock()
		mr, ok := ms.responses[ep.Name]
		ms.mu.RUnlock()

		if ok && mr.handler != nil {
			mr.handler(w, r)
			return
		}

		statusCode := http.StatusOK
		var data any

		if ok {
			if mr.statusCode != nil {
				statusCode = *mr.statusCode
			}
			data = mr.data
		}

		if statusCode >= 300 {
			ms.writeErrorResponse(w, ep, statusCode)
			return
		}

		if data != nil {
			ms.writeJSONResponse(w, ep, statusCode, data)
			return
		}

		switch ep.Backend {
		case cav.BackendInfrapi:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cloud-Avenue-Mock", "true")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"jobId":"87ab1934-0146-4fb0-80bc-815fea03214d","message":"Job created successfully"}`)) //nolint:errcheck

		case cav.BackendVMware:
			w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
			w.Header().Set("X-Cloud-Avenue-Mock", "true")
			w.WriteHeader(http.StatusAccepted)
		}
	}
}
