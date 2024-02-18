package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

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
