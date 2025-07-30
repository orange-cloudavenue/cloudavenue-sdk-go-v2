package edgegateway

import (
	"github.com/orange-cloudavenue/common-go/urn"
)

type (
	// * apiResponse
	apiResponseT0s []apiResponseT0

	apiResponseT0 struct {
		Type       string                  `json:"type" fake:"tier-0-vrf"`
		Name       string                  `json:"name" fake:"{resource_name:t0}"`
		Properties apiResponseT0Properties `json:"properties,omitempty"`
		Children   []apiResponseT0Children `json:"children,omitempty" fakesize:"1"`
	}

	apiResponseT0Properties struct {
		ClassOfService string `json:"classOfService,omitempty" fake:"SHARED_STANDARD"`
	}

	apiResponseT0Children struct {
		Type       string `json:"type" fake:"edge-gateway"`
		Name       string `json:"name" fake:"{resource_name:edgegateway}"`
		Properties struct {
			RateLimit int    `json:"rateLimit,omitempty" fake:"5"`
			EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"` // The UUID of the edge gateway
		} `json:"properties,omitempty"`
	}

	// * Params

	ParamsGetT0 struct {
		T0Name          string
		EdgegatewayID   string
		EdgegatewayName string
	}

	//* Model

	// ModelT0s represents the model for T0 routers.
	ModelT0s struct {
		Count int
		T0s   []ModelT0
	}

	// ModelT0 represents a T0 router model.
	// It contains the name, class of service, bandwidth, maximum edge gateways,
	ModelT0 struct {
		// Name defines the name of the T0 router.
		Name string `documentation:"Name of the T0 router"`

		// ClassOfService defines the class of service for the T0 router.
		ClassOfService string `documentation:"Class of service for the T0 router."`

		// Bandwidth defines the bandwidth for the T0 router.
		Bandwidth ModelT0Bandwidth `documentation:"Bandwidth for the T0 router in Mbps."`

		// MaxEdgeGateways defines the maximum number of edge gateways for the T0 router.
		// This is a limit imposed by the Class of Service.
		MaxEdgeGateways int `documentation:"Maximum number of edge gateways for the T0 router."`

		// EdgeGateways contains the list of edge gateways associated with the T0 router.
		EdgeGateways []ModelT0EdgeGateway `documentation:"List of edge gateways associated with the T0 router."`
	}

	// ModelT0Bandwidth represents the bandwidth model for a T0 router.
	ModelT0Bandwidth struct {
		// Capacity defines the total bandwidth capacity in Mbps.
		// This is the maximum bandwidth that can be allocated to the T0 router.
		// It is a limit imposed by the Class of Service.
		Capacity int `documentation:"Total bandwidth capacity for the T0 router in Mbps. This is the maximum bandwidth that can be allocated to the T0 router. It is a limit imposed by the Class of Service."`

		// Provisioned defines the amount of bandwidth that has been provisioned in the each edge gateway.
		// This is the total bandwidth that has been allocated to the T0 router across all edge gateways.
		// It is the sum of the bandwidth allocated to each edge gateway.
		Provisioned int `documentation:"Total bandwidth provisioned for the T0 router across all edge gateways in Mbps."`

		// Remaining defines the remaining bandwidth that can be allocated to the edge gateways.
		Remaining int `documentation:"Remaining bandwidth that can be allocated to the edge gateways in Mbps. This is calculated as Capacity - Provisioned."`

		// AllowedBandwidthValues returns the allowed bandwidth values for the T0 router.
		// It's used to determine the available bandwidth options for the new edge gateway.
		// It returns a slice of integers representing the allowed bandwidth values in Mbps.
		// If values are empty, no bandwidth remaining to allocate.
		AllowedBandwidthValues []int `documentation:"Allowed bandwidth values for the T0 router in Mbps. This is used to determine the available bandwidth options for the new edge gateway. If empty, no bandwidth remaining to allocate."`

		// AllowUnlimited indicates if unlimited bandwidth is allowed for the T0 router.
		// This is true if the T0 router is DEDICATED.
		AllowUnlimited bool `documentation:"Indicates if unlimited bandwidth is allowed for the T0 router. This is true if the T0 router is DEDICATED."`
	}

	// ModelT0EdgeGateway represents an edge gateway associated with a T0 router.
	ModelT0EdgeGateway struct {
		// ID defines the unique identifier of the edge gateway.
		// It is a URN that uniquely identifies the edge gateway.
		ID string `documentation:"Unique identifier of the edge gateway. It is a URN that uniquely identifies the edge gateway."`

		// Name defines the name of the edge gateway.
		Name string `documentation:"Name of the edge gateway."`

		// Bandwidth defines the bandwidth allocated to the edge gateway.
		// The value is in Mbps.
		Bandwidth int `documentation:"Bandwidth allocated to the edge gateway in Mbps."`

		// AllowedBandwidthValues returns the allowed bandwidth values for the edge gateway.
		AllowedBandwidthValues []int `documentation:"Allowed bandwidth values for the edge gateway in Mbps."`
	}
)

func (t0s apiResponseT0s) toModel() *ModelT0s {
	var modelT0s ModelT0s
	for _, t0 := range t0s {
		modelT0 := ModelT0{
			Name:           t0.Name,
			ClassOfService: t0.Properties.ClassOfService,
			Bandwidth: ModelT0Bandwidth{
				Capacity: func() int {
					if cof, ok := classOfServices[t0.Properties.ClassOfService]; ok {
						return cof.MaxBandwidth
					}
					return 0
				}(),
				Provisioned: func() int {
					if t0.Children == nil {
						return 0
					}

					var totalProvisioned int
					for _, child := range t0.Children {
						if child.Properties.RateLimit > 0 {
							totalProvisioned += child.Properties.RateLimit
						}
					}
					return totalProvisioned
				}(),
				Remaining:              0,   // Remaining is calculated later
				AllowedBandwidthValues: nil, // AllowedBandwidthValues is calculated later
				AllowUnlimited: func() bool {
					if cof, ok := classOfServices[t0.Properties.ClassOfService]; ok {
						return cof.AllowUnlimited
					}
					return false
				}(),
			},
			MaxEdgeGateways: func() int {
				if cof, ok := classOfServices[t0.Properties.ClassOfService]; ok {
					return cof.MaxEdgeGateways
				}
				return 0
			}(),
			EdgeGateways: func() []ModelT0EdgeGateway {
				var edgeGateways []ModelT0EdgeGateway

				for _, child := range t0.Children {
					if child.Type == "edge-gateway" {
						edgeGateway := ModelT0EdgeGateway{
							ID:                     urn.Normalize(urn.EdgeGateway, child.Properties.EdgeUUID).String(),
							Name:                   child.Name,
							Bandwidth:              child.Properties.RateLimit,
							AllowedBandwidthValues: nil, // This will be calculated later based on the remaining bandwidth
						}
						edgeGateways = append(edgeGateways, edgeGateway)
					}
				}
				return edgeGateways
			}(),
		}
		modelT0s.T0s = append(modelT0s.T0s, modelT0)
	}

	// Calculate the remaining bandwidth and allowed bandwidth values for each T0 router
	for i, t0 := range modelT0s.T0s {
		t0.Bandwidth.Remaining = t0.Bandwidth.Capacity - t0.Bandwidth.Provisioned

		// Calculate the allowed bandwidth values for the T0 router
		if cof, ok := classOfServices[t0.ClassOfService]; ok {
			// Use t0.Bandwidth.Remaining to ensure the values are within the allowed range
			if t0.Bandwidth.Remaining > 0 {
				for _, bw := range cof.MaxEdgeGatewayBandwidth {
					if bw <= t0.Bandwidth.Remaining {
						t0.Bandwidth.AllowedBandwidthValues = append(t0.Bandwidth.AllowedBandwidthValues, bw)
					}
				}
			}
		}
		modelT0s.T0s[i] = t0
	}

	// Calculate the allowed bandwidth values for each edge gateway
	for i, t0 := range modelT0s.T0s {
		for j, edgeGateway := range t0.EdgeGateways {
			// Use t0.Bandwidth.Remaining to ensure the values are within the allowed range
			if t0.Bandwidth.Remaining > 0 {
				for _, bw := range classOfServices[t0.ClassOfService].MaxEdgeGatewayBandwidth {
					if bw <= t0.Bandwidth.Remaining {
						edgeGateway.AllowedBandwidthValues = append(edgeGateway.AllowedBandwidthValues, bw)
					}
				}
			}
			modelT0s.T0s[i].EdgeGateways[j] = edgeGateway
		}
	}

	modelT0s.Count = len(modelT0s.T0s)
	return &modelT0s
}
