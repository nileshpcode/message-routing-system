package sqlite

import (
	"database/sql"
	"io/ioutil"
	"log"
)

// db wrapper
type DB struct {
	*sql.DB // sql db
}

// db service.
type DBSvc struct {
	Dbo    *DB
	logger *log.Logger
}

// Open opens a new db connection.
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// create test schema
func (svc *DBSvc) CreateTestSchema() error {
	if err := svc.DropAndCreateTable(GatewayTableFilePath); err != nil {
		return err
	}
	if err := svc.DropAndCreateTable(RouteTableFilePath); err != nil {
		return err
	}

	return nil
}

// create table
func (svc *DBSvc) DropAndCreateTable(tblFilePath string) error {
	tbl, err := ioutil.ReadFile(tblFilePath)
	if err != nil {
		return err
	}

	_, err = svc.Dbo.Exec(string(tbl))
	return err
}
