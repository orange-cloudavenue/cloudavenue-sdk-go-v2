package edgegateway

// * Models

type (
	// ModelEdgeGateways represents a list of edge gateways.
	ModelEdgeGateways struct {
		EdgeGateways []ModelEdgeGateway
	}
	// ModelEdgeGateway represents the model of an edge gateway.
	ModelEdgeGateway struct {
		ID string

		// Name of edge gateway
		Name string

		// Description of edge gateway
		Description string

		// OwnerRef defines VDC or VDC Group that this network belongs to.
		OwnerRef *ModelObjectReference

		// UplinkT0 defines the T0 router name that this edge gateway is connected to.
		UplinkT0 *ModelObjectReference

		// Services is the list of network services
		// that are available on the edge gateway
		// Services ModelNetworkServicesSvcs
	}

	// Bandwidth in Mbps.
	ModelBandwidth int
)

// * Functions Parameters

type (
	ParamsEdgeGateway struct {
		ID   string `validate:"required_if_null=Name,omitempty,urn=edgeGateway" fake:"{urn:edgeGateway}"`
		Name string `validate:"required_if_null=ID,omitempty,edgegateway_name" fake:"{edgegateway_name}"`
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
		T0Name string `fake:"{t0_name}"`
	}
)

// * Request / Response API

type (
	apiRequestEdgeGateway struct {
		T0Name string `json:"tier0VrfId" fake:"{t0_name}"`
	}

	apiResponseEdgegateways struct {
		Values []apiResponseEdgegateway `json:"values,omitempty" fakesize:"1"` // List of edge gateways.
	}
	apiResponseEdgegateway struct {
		ID          string `json:"id" fake:"{urn:edgeGateway}"`    // The ID of the edge gateway.
		Name        string `json:"name" fake:"{edgegateway_name}"` // The name of the edge gateway.
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
			UplinkName string `json:"uplinkName" fake:"{t0_name}"` // The name of the uplink.
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
		Name                string `json:"name,omitempty" fake:"{edgegateway_name}"`          // EdgeGateway name.
		VDCID               string `json:"vdc,omitempty" fake:"{urn:vdc}"`                    // VDC Reference or ID
		VDCName             string `json:"orgVdcName,omitempty" fake:"{word}"`                // VDC name
		NumberOfExtNetworks int    `json:"numberOfExtNetworks,omitempty" fake:"{number:1,5}"` // Number of external networks connected to the edgeGateway.
		NumberOfOrgNetworks int    `json:"numberOfOrgNetworks,omitempty" fake:"{number:1,5}"` // Number of org VDC networks connected to the edgeGateway
		IsBusy              bool   `json:"isBusy,omitempty" fake:"{bool}"`                    // True if this Edge Gateway is busy.
		GatewayStatus       string `json:"gatewayStatus,omitempty" fake:"{word}"`             // Status of the edgeGateway
	}

	apiResponseBandwidth struct {
		// bandwidth in Mbps. This value is not returned by the API
		// It will be set by the middleware and the value are extracted
		// from the egress and ingress profiles names. (e.g. "qosgw005mbps" = 5mbps)
		// nil value means no bandwidth limit.
		Bandwidth *int `json:"-" fake:"-"`

		// In cloudavenue egress and ingress profiles are always the same.
		// Use only egressProfile to get the bandwidth.
		EgressProfile struct {
			Name string `json:"name" fake:"{randomstring:[qosgw005mbps,qosgw025mbps,qosgw050mbps,qosgw075mbps,qosgw100mbps,qosgw150mbps,qosgw200mbps,qosgw250mbps,qosgw300mbps]}"`
		} `json:"egressProfile"`
		IngressProfile struct{} `json:"-"`
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

// toModel converts the apiResponseBandwidth to ModelBandwidth.
func (api *apiResponseBandwidth) toModel() *ModelBandwidth {
	if api == nil {
		return nil
	}

	if api.Bandwidth == nil {
		return nil // No bandwidth limit
	}

	return (*ModelBandwidth)(api.Bandwidth)
}
