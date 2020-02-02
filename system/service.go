package system

type Service interface {
	CreateGateway(gateway *Gateway) error
	GetGateway(id int64) (*Gateway, error)
	GetRoute(id int64) (*Route, error)
	CreateRoute(route *Route) error
	QueryGateway(params map[GatewayQueryParam]interface{}) ([]*Gateway, error)
	QueryRoute(params map[RouteQueryParam]interface{}) ([]*Route, map[string]*Route, error)
}
