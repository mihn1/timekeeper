package datatypes

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DateOnly represents a date without a time component
type DateOnly struct {
	time.Time
}

// NewDateOnly creates a DateOnly from a time.Time
func NewDateOnly(t time.Time) DateOnly {
	// Normalize the time to remove time part
	normalized := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return DateOnly{Time: normalized}
}

func NewDateOnlyFromStr(dateStr string) (DateOnly, error) {
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return DateOnly{}, err
	}
	return NewDateOnly(parsedTime), nil
}

// Scan implements the sql.Scanner interface
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		*d = DateOnly{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = NewDateOnly(v)
	case string:
		parsedTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		*d = NewDateOnly(parsedTime)
	case []byte:
		parsedTime, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		*d = NewDateOnly(parsedTime)
	default:
		return fmt.Errorf("cannot scan type %T into DateOnly", value)
	}
	return nil
}

// Value implements the driver.Valuer interface
func (d DateOnly) Value() (driver.Value, error) {
	return d.Time, nil
}

// String returns the date in YYYY-MM-DD format
func (d DateOnly) String() string {
	return d.Format("2006-01-02")
}
