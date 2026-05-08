package common
type ProxyStats struct {
	AvgRewriteDelayNs int
	AvgModifyResponseAssetDeliveryDelayNs int
	AvgModifyResponseCdnDelayNs int
}

type SocketMessage struct {
	Type string
	Error string //an error of "" should be interpreted as success
	DataB []byte //generic byte array that may be used for data
	Stats ProxyStats
}