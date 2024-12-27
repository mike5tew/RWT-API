// package for processing the site information
package tools

import (
	"time"
)

func ParseMySQLDateTime(datetimeStr string) (time.Time, error) {
	if len(datetimeStr) == 10 { // Check if the string is a date without time
		datetimeStr += " 00:00:00" // Append default time
	}
	return time.Parse("2006-01-02 15:04:05", datetimeStr)
}
