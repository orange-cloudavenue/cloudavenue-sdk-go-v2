package edgegateway

type (
	ModelEdgeGatewayBandwidth ModelT0EdgeGateway

	apiRequestBandwidth struct {
		Bandwidth int `json:"rateLimit" fake:"5"`
	}
)
