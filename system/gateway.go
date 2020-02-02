package system

type Gateway struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IpAddresses Ips    `json:"ip_addresses"`
}

type GatewayQueryParam string

const (
	GatewayQueryParamName = GatewayQueryParam("name")
)
