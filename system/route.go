package system

type Route struct {
	ID      int64    `json:"id"`
	Prefix  string   `json:"prefix"`
	Gateway *Gateway `json:"gateway"`
}

type RouteQueryParam string

const (
	RouteQueryParamPrefix = RouteQueryParam("prefix")
)
