package cav

import (
	"context"
	"fmt"

	"resty.dev/v3"
)

// NewRequest creates a new request using the resty client.
func (c *client) NewRequest(ctx context.Context, client SubClientName, opts ...RequestOption) (req *resty.Request, err error) {
	// Retrieve the subclient based on the provided client name.
	// This method identifies the subclient and returns it.
	sc, err := c.identifyClient(ctx, client)
	if err != nil {
		return nil, err
	}

	// Inject the client name into the context for retrieval in the other methods.
	ctxv := context.WithValue(ctx, contextKeyClientName, client)

	// Create and populate the request options.
	reqOpts, err := newRequestOptions(opts...)
	if err != nil {
		return nil, err
	}

	// TODO(azrod) thinking about removing the error from the new HTTPClient because
	// it is not really useful, we can just return the request with the context. No error
	// should be returned here, because the client should always be able to create a new HTTP.

	// Setup the HTTP client for the request.
	// This client is used to send the request and handle the response.
	hC, err := sc.NewHTTPClient(ctxv)
	if err != nil {
		return nil, err
	}

	// If JobOpts are provided, we need to create a request with job middleware.
	// This is used to handle job responses and status checks.
	if reqOpts.JobOpts != nil {
		sCJob, ok := sc.(jobsInterface)
		if !ok {
			// If the subclient does not implement the jobs.Client interface,
			// we cannot create a job request.
			// Return an error indicating that the client does not support job options.
			// This is a design choice to ensure that only clients that support jobs can use this feature.
			// If you need to handle jobs, ensure that the client implements the jobs.Client interface.
			return nil, fmt.Errorf("client %s does not support job options", client)
		}

		// Create a new HTTP client specifically for job requests.
		// This is necessary because the initial client (hc) have a specific middleware defined below.
		// If the hc client has used in NewJobMiddleware, it will create an infinite loop.
		// So we create a new client for job requests.
		hCForJob, err := sc.NewHTTPClient(ctxv)
		if err != nil {
			return nil, err
		}

		// If the request is for a job, set the job middleware.
		hC.SetResponseMiddlewares(
			resty.AutoParseResponseMiddleware,
			newJobMiddleware(hCForJob, sCJob, reqOpts.JobOpts),
		)
	}

	// Create a new request with the context and options.
	hR := hC.NewRequest().
		SetContext(ctxv)

	return hR, nil
}
