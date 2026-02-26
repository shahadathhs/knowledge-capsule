package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Tags stores a string slice as JSON in the database.
type Tags []string

// Value implements driver.Valuer for GORM.
func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return "[]", nil
	}
	return json.Marshal(t)
}

// Scan implements sql.Scanner for GORM.
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("tags: invalid type for scan")
	}
	return json.Unmarshal(bytes, t)
}
