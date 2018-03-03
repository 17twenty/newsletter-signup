package main

import (
	"database/sql/driver"
	"time"
)

// SqlLiteDate is a nice time function ... mainly as SQLite is shit at dates
type SqlLiteDate time.Time

// Scan is used to extract the value from the database (comes out as []uint8)
func (b *SqlLiteDate) Scan(value interface{}) error {
	t, _ := time.Parse("2006-01-02", string(value.([]uint8)))
	*b = SqlLiteDate(t)
	return nil
}

// Value is used to the date tiem as a string so it goes into the database
func (b SqlLiteDate) Value() (driver.Value, error) {
	// log.Println("Value:", time.Time(b).Format("2006-01-02"))
	return time.Time(b).Format("2006-01-02"), nil
}

func (b SqlLiteDate) String() string {
	return time.Time(b).Format("2006-01-02")
}
