package date

import "time"

func NumberOfDaysInMonth(m time.Month, year int) int64 {
	return int64(time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day())
}

func NumberOfDaysInTime(date time.Time) int {
	year, month, _ := date.Date()
	return time.Date(year, month+1, 0, 0, 0, 0, 0, date.Location()).Day()
}
