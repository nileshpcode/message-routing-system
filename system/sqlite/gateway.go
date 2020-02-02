package sqlite

import (
	"errors"
	"fmt"
	"question/system"
	"strings"
)

const GatewayTableFilePath = "system/sqlite/gateway.tbl"

func (db *DBSvc) CreateGateway(gateway *system.Gateway) error {

	result, err := db.Dbo.Exec(fmt.Sprintf(`INSERT INTO gateway (name, ip_addresses) VALUES (?, ?)`), gateway.Name, gateway.IpAddresses)
	if err != nil {
		return err
	}

	gateway.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (db *DBSvc) GetGateway(id int64) (*system.Gateway, error) {
	q := fmt.Sprintf(`select name, ip_addresses FROM gateway WHERE id = ?`)

	rows, err := db.Dbo.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var name string
	var ips system.Ips

	for rows.Next() {
		if err := rows.Scan(&name, &ips); err != nil {
			return nil, err
		}
		return &system.Gateway{
			ID:          id,
			Name:        name,
			IpAddresses: ips,
		}, nil
	}

	return nil, errors.New("gateway id not found")
}

func (db *DBSvc) QueryGateway(params map[system.GatewayQueryParam]interface{}) ([]*system.Gateway, error) {
	query := `SELECT id, name, ip_addresses FROM gateway WHERE %s`
	wheres := []string{`1=1`}

	args := make([]interface{},0)
	if params != nil {
		if param, ok := params[system.GatewayQueryParamName]; ok {
			args = append(args, param)
			wheres = append(wheres, fmt.Sprintf(`name = ?`))
		}
	}

	query = fmt.Sprintf(query, strings.Join(wheres, " AND "))

	rows, err := db.Dbo.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	gateways := make([]*system.Gateway, 0)

	var id int64
	var name string
	var ips system.Ips

	for rows.Next() {
		if err := rows.Scan(&id, &name, &ips); err != nil {
			return nil, err
		}
		g := &system.Gateway{
			ID:          id,
			Name:        name,
			IpAddresses: ips,
		}
		gateways = append(gateways, g)
	}

	return gateways, nil
}
