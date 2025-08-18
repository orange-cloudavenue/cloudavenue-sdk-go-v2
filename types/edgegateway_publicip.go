package types

type (
	ModelEdgeGatewayPublicIPs struct {
		EdgegatewayID   string `documentation:"ID of the edge gateway"`
		EdgegatewayName string `documentation:"Name of the edge gateway"`

		PublicIPs []ModelEdgeGatewayPublicIP `documentation:"List of public IPs"`
	}

	ModelEdgeGatewayPublicIP struct {
		EdgegatewayID   string `documentation:"ID of the edge gateway"`
		EdgegatewayName string `documentation:"Name of the edge gateway"`

		ModelEdgeGatewayServicesPublicIP
	}

	ParamsGetEdgeGatewayPublicIP struct {
		ID   string `fake:"{urn:edgegateway}"`
		Name string `fake:"{resource_name:edgegateway}"`

		IP string `validate:"ipv4" fake:"{IPv4Address}"`
	}

	ParamsDeleteEdgeGatewayPublicIP struct {
		IP string `validate:"ipv4" fake:"{IPv4Address}"`
	}
)
