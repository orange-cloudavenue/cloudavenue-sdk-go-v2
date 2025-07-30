package edgegateway

// * Models

type (
	// ModelEdgeGateways represents a list of edge gateways.
	ModelEdgeGateways struct {
		EdgeGateways []ModelEdgeGateway
	}
	// ModelEdgeGateway represents the model of an edge gateway.
	ModelEdgeGateway struct {
		ID string `documentation:"ID of the edge gateway"`

		// Name of edge gateway
		Name string `documentation:"Name of the edge gateway"`

		// Description of edge gateway
		Description string `documentation:"Description of the edge gateway"`

		// OwnerRef defines VDC or VDC Group that this network belongs to.
		OwnerRef *ModelObjectReference `documentation:"VDC or VDC Group that this edge gateway belongs to"`

		// UplinkT0 defines the T0 router name that this edge gateway is connected to.
		UplinkT0 *ModelObjectReference `documentation:"T0 router name that this edge gateway is connected to"`
	}

	// Bandwidth in Mbps.
	ModelBandwidth int
)

// * Functions Parameters

type (
	ParamsEdgeGateway struct {
		ID   string `validate:"required_if_null=Name,omitempty,urn=edgegateway" fake:"{urn:edgegateway}"`
		Name string `validate:"required_if_null=ID,omitempty,resource_name=edgegateway" fake:"{resource_name:edgegateway}"`
	}

	ParamsCreateEdgeGateway struct {
		// OwnerType is the type of the owner of the edge gateway.
		// It can be either "vdc" or "vdcgroup".
		OwnerType string `fake:"{randomstring:[vdc]}"`

		// OwnerName is the VDC or VDC Group that this edge gateway belongs to.
		OwnerName string `fake:"{vdc_name}"`

		// Name is the name of the T0 router that this edge gateway will be connected to.
		// If not provided and only if one T0 router is available,
		// the first T0 router will be used.
		T0Name string `fake:"{resource_name:t0}"`

		// Bandwidth is the bandwidth limit in Mbps.
		// If not provided, default bandwidth will be used (5 Mbps).
		// You can get bandwidth available values for the new edge gateway
		// by calling ListT0().
		// Unlimited bandwidth is allowed if the T0 is DEDICATED.
		Bandwidth int `fake:"5"`
	}

	ParamsUpdateEdgeGateway struct {
		// ID is the ID of the edge gateway to update.
		ID string `fake:"{urn:edgegateway}"`

		// Name is the new name of the edge gateway.
		Name string `fake:"{resource_name:edgegateway}"`

		// Bandwidth is the new bandwidth limit in Mbps.
		Bandwidth int `fake:"5"`
	}
)

// * Request / Response API

type (
	apiRequestEdgeGateway struct {
		T0Name string `json:"tier0VrfId" fake:"{resource_name:t0}"`
	}

	apiResponseEdgegateways struct {
		Values []apiResponseEdgegateway `json:"values,omitempty" fakesize:"1"` // List of edge gateways.
	}
	apiResponseEdgegateway struct {
		ID          string `json:"id" fake:"{urn:edgegateway}"`             // The ID of the edge gateway.
		Name        string `json:"name" fake:"{resource_name:edgegateway}"` // The name of the edge gateway.
		Description string `json:"description" fake:"{sentence}"`

		EdgeGatewayUplinks []struct {
			Connected bool `json:"connected" fake:"true"` // Indicates if the uplink is connected.
			Dedicated bool `json:"dedicated"`
			Subnets   struct {
				Values []struct {
					DNSServer1 string `json:"dnsServer1" fake:"{ipv4address}"`
					DNSServer2 string `json:"dnsServer2" fake:"{ipv4address}"`
					DNSSuffix  string `json:"dnsSuffix" fake:"{domainname}"`
					Enabled    bool   `json:"enabled" fake:"{bool}"` // Indicates if the subnet is enabled.
					Gateway    string `json:"gateway" fake:"{ipv4address}"`
					IPRanges   struct {
						Values []struct {
							EndAddress   string `json:"endAddress" fake:"{ipv4address}"`
							StartAddress string `json:"startAddress" fake:"{ipv4address}"`
						} `json:"values" fakesize:"1"`
					} `json:"ipRanges"`
					PrefixLength int64  `json:"prefixLength" fake:"{number:24,32}"` // The prefix length of the subnet.
					PrimaryIP    string `json:"primaryIp" fake:"{ipv4address}"`
					TotalIPCount int64  `json:"totalIpCount" fake:"{number:5,10}"` // The total number of IP addresses in the subnet.
					UsedIPCount  int64  `json:"usedIpCount" fake:"{number:6,8}"`   // The number of used IP addresses in the subnet.
				} `json:"values" fakesize:"1"`
			} `json:"subnets" fakesize:"1"`
			UplinkID   string `json:"uplinkId" fake:"{urn:network}"`
			UplinkName string `json:"uplinkName" fake:"{resource_name:t0}"` // The name of the uplink.
		} `json:"edgeGatewayUplinks" fakesize:"1"`

		OrgVDC *ModelObjectReference `json:"orgVdc"`

		// OwnerRef contains information about the owner of the edge gateway (VDC Or VDCGroup)
		OwnerRef *ModelObjectReference `json:"ownerRef"`

		// OrgVdcNetworkCount holds the number of Org VDC networks connected to the gateway.
		OrgVDCNetworkCount int64 `json:"orgVdcNetworkCount" fake:"{number:1,10}"`
	}

	apiResponseQueryEdgeGateway struct {
		Record []apiResponseQueryEdgeGatewayRecord `json:"record,omitempty" fakesize:"1"` // List of edge gateways.
	}

	apiResponseQueryEdgeGatewayRecord struct {
		ID                  string `json:"-"`                                                 // The ID of the entity.
		HREF                string `json:"href,omitempty" fake:"{href_uuid}"`                 // The URI of the entity.
		Type                string `json:"type,omitempty"`                                    // The MIME type of the entity.
		Name                string `json:"name,omitempty" fake:"{resource_name:edgegateway}"` // EdgeGateway name.
		VDCID               string `json:"vdc,omitempty" fake:"{urn:vdc}"`                    // VDC Reference or ID
		VDCName             string `json:"orgVdcName,omitempty" fake:"{word}"`                // VDC name
		NumberOfExtNetworks int    `json:"numberOfExtNetworks,omitempty" fake:"{number:1,5}"` // Number of external networks connected to the edgeGateway.
		NumberOfOrgNetworks int    `json:"numberOfOrgNetworks,omitempty" fake:"{number:1,5}"` // Number of org VDC networks connected to the edgeGateway
		IsBusy              bool   `json:"isBusy,omitempty" fake:"{bool}"`                    // True if this Edge Gateway is busy.
		GatewayStatus       string `json:"gatewayStatus,omitempty" fake:"{word}"`             // Status of the edgeGateway
	}
)

// toModel converts the apiResponseEdgegateways to ModelEdgeGateways.
func (api *apiResponseEdgegateways) toModel() *ModelEdgeGateways {
	if api == nil {
		return nil
	}

	model := &ModelEdgeGateways{
		EdgeGateways: make([]ModelEdgeGateway, 0, len(api.Values)),
	}

	for _, v := range api.Values {
		model.EdgeGateways = append(model.EdgeGateways, *v.toModel())
	}

	return model
}

// toModel converts the edgegatewayAPIResponse to ModelEdgeGateway.
func (api *apiResponseEdgegateway) toModel() *ModelEdgeGateway {
	if api == nil {
		return nil
	}

	return &ModelEdgeGateway{
		ID:          api.ID,
		Name:        api.Name,
		Description: api.Description,
		OwnerRef:    api.OwnerRef,
		UplinkT0: func() *ModelObjectReference {
			if len(api.EdgeGatewayUplinks) > 0 {
				return &ModelObjectReference{
					ID:   api.EdgeGatewayUplinks[0].UplinkID,
					Name: api.EdgeGatewayUplinks[0].UplinkName,
				}
			}
			return nil
		}(),
	}
}
