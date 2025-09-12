package utils

import "time"

func OffsetDate(startTime time.Time, days int) time.Time {
	y, m, d := startTime.Date()
	utcMidnight := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	expiryDate := utcMidnight.AddDate(0, 0, days)
	return expiryDate
}
