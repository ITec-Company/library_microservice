package utils

import (
	"database/sql/driver"
	"time"
)

type AnyTime struct {
	time.Time
}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
func (a AnyTime) Value() (driver.Value, error) {
	return driver.Value(a.Time), nil
}
