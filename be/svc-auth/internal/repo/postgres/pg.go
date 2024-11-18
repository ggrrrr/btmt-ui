package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

const createPasswdTable string = `CREATE TABLE IF NOT EXISTS %s (
	subject TEXT  not null,
	passwd TEXT  not null,
	"status" TEXT  not null,
	system_roles TEXT[] not null,
	tenant_roles JSONB  not null,
	created_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(subject)
)`

const createHistoryTable string = `CREATE TABLE IF NOT EXISTS %s (
	"id" TEXT not null,
	"subject" TEXT not null,
	"method" TEXT not null,
	"device" JSONB not null,
	"created_at" TIMESTAMP DEFAULT NOW(),
	UNIQUE(id)
)`

type (
	tRoles   map[string][]string
	authRepo struct {
		passwdTable  string
		historyTable string
		db           *sql.DB
	}
)

var _ (ddd.AuthPasswdRepo) = (*authRepo)(nil)
var _ (ddd.AuthHistoryRepo) = (*authRepo)(nil)

func (r *authRepo) table(sql string, tableName string) string {
	return fmt.Sprintf(sql, tableName)
}

func Init(db *sql.DB) (*authRepo, error) {
	r := &authRepo{
		passwdTable:  "auth_passwd",
		historyTable: "auth_history",
		db:           db,
	}
	err := create(r.db, r.table(createPasswdTable, r.passwdTable))
	if err != nil {
		return nil, err
	}
	err = create(r.db, r.table(createHistoryTable, r.historyTable))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func create(db *sql.DB, sql string) error {
	logger.Info().Str("sql", sql).Msg("create table")
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

type device struct {
	app.Device
}

func (v *device) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &v.Device)
}

func (a device) Value() (driver.Value, error) {
	return json.Marshal(a.Device)
}

func (v *tRoles) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &v)
}

func (a tRoles) Value() (driver.Value, error) {
	return json.Marshal(a)
}
