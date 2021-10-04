package dates

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// CapuchinFormat is the format of the date that Capuchin understands
	// it's zero padded for months and days, and uses a four figure year,
	// e.g. 31122021
	CapuchinFormat = "%02d%02d%d"
)

// ParseDate converts whatever date format is given to 01012012 format
func ParseDate(date, format string) (int, error) {
	t, err := time.Parse(format, date)
	if err != nil {
		return 0, fmt.Errorf("error parsing date with layout %s: %v", format, err)
	}

	formatted := fmt.Sprintf(CapuchinFormat, t.Day(), t.Month(), t.Year())

	return strconv.Atoi(formatted)
}
