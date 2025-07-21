package edgegateway

type (
	// ModelObjectReference represents a reference to an object in the API.
	// It contains the ID and name of the object.
	ModelObjectReference struct {
		ID   string `json:"id" fake:"{urn:vdc}"`
		Name string `json:"name"`
	}
)
