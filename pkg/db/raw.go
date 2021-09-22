// Date: 2021/9/20

// Package db
package db

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.
type RawMessage []byte

// Scan implements the Scanner interface.
func (m *RawMessage) Scan(value interface{}) error {
	if value == nil {
		*m = append((*m)[0:0], []byte("null")...)
	}
	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("%v not byte", value)
	}
	*m = append((*m)[0:0], v...)
	return nil
}

// Value implements the driver Valuer interface.
func (m RawMessage) Value() (driver.Value, error) {
	if len(m) == 0 {
		return []byte("null"), nil
	}
	return []byte(m), nil
}

// MarshalJSON returns m as the JSON encoding of m.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}
