package cav

type (
	endpointRegistryOptions struct {
		// API is the API name.
		api API
		// Version is the API version.
		version Version
	}

	EndpointRegistryOptions func(*endpointRegistryOptions)
)

// WithExtraProperties allows adding extra properties to the endpoint registry options.
func WithExtraProperties(api API, version Version) EndpointRegistryOptions {
	return func(opts *endpointRegistryOptions) {
		opts.api = api
		opts.version = version
	}
}
