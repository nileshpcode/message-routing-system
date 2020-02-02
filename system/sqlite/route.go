package sqlite

import (
	"errors"
	"fmt"
	"question/system"
	"strings"
)

const RouteTableFilePath = "system/sqlite/route.tbl"

func (db *DBSvc) CreateRoute(route *system.Route) error {

	result, err := db.Dbo.Exec(fmt.Sprintf(`INSERT INTO route (prefix, gateway_id) VALUES (?, ?)`), route.Prefix, route.Gateway.ID)
	if err != nil {
		return err
	}
	route.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (db *DBSvc) GetRoute(id int64) (*system.Route, error) {
	q := fmt.Sprintf(`
		SELECT 
			route.prefix,
			route.gateway_id,
			gateway.name,
			gateway.ip_addresses
		FROM 
			route 
			JOIN gateway ON gateway.id = route.gateway_id
		WHERE 
			route.id = ?`)

	rows, err := db.Dbo.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prefix string
	var gatewayId int64
	var gatewayName string
	var gatewayIps system.Ips

	for rows.Next() {
		if err := rows.Scan(&prefix, &gatewayId, &gatewayName, &gatewayIps); err != nil {
			return nil, err
		}

		gateway := &system.Gateway{
			ID:          gatewayId,
			Name:        gatewayName,
			IpAddresses: gatewayIps,
		}

		return &system.Route{
			ID:      id,
			Prefix:  prefix,
			Gateway: gateway,
		}, nil

	}

	return nil, errors.New("route id not found")
}

func (db *DBSvc) QueryRoute(params map[system.RouteQueryParam]interface{}) ([]*system.Route, map[string]*system.Route, error) {
	query := `
		SELECT 
			route.id,
			route.prefix,
			route.gateway_id,
			gateway.name,
			gateway.ip_addresses
		FROM 
			route 
			JOIN gateway ON gateway.id = route.gateway_id
		WHERE 
			%s
			`
	wheres := []string{`1=1`}
	args := make([]interface{}, 0)
	if params != nil {
		if param, ok := params[system.RouteQueryParamPrefix]; ok {
			args = append(args, param)
			wheres = append(wheres, fmt.Sprintf(`route.prefix = ?`))
		}
	}
	query = fmt.Sprintf(query, strings.Join(wheres, " AND "))

	rows, err := db.Dbo.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	routes := make([]*system.Route, 0)

	var id int64
	var prefix string
	var gatewayID int64
	var gatewayName string
	var gatewayIps system.Ips

	prefixMap := make(map[string]*system.Route)
	for rows.Next() {
		if err := rows.Scan(&id, &prefix, &gatewayID, &gatewayName, &gatewayIps); err != nil {
			return nil, nil, err
		}
		g := &system.Route{
			ID:     id,
			Prefix: prefix,
			Gateway: &system.Gateway{
				ID:          gatewayID,
				Name:        gatewayName,
				IpAddresses: gatewayIps,
			},
		}
		prefixMap[prefix] = g

		routes = append(routes, g)
	}

	return routes, prefixMap, nil
}
