package system

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Ip string
type Ips []Ip

// Value implements the driver.Valuer interface
func (ips Ips) Value() (driver.Value, error) {
	ipsStrings := make([]string, 0)
	for _, ip := range ips {
		ipsStrings = append(ipsStrings, string(ip))
	}

	return strings.Join(ipsStrings, ","), nil
}

// Scan implements the sql.Scanner interface
func (ips *Ips) Scan(src interface{}) error {

	if v, ok := src.(string); ok {
		ipsStringsArr := strings.Split(v, ",")
		aux := Ips{}

		for _, ip := range ipsStringsArr {
			aux = append(aux, Ip(ip))
		}
		*ips = aux

		return nil
	}

	return fmt.Errorf("unable to scan %v as string", src)
}
