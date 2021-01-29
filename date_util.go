package webutility

import (
	"fmt"
	"strings"
	"time"
)

const (
	YYYYMMDD_sl = "2006/01/02"
	YYYYMMDD_ds = "2006-01-02"
	YYYYMMDD_dt = "2006.01.02."

	DDMMYYYY_sl = "02/01/2006"
	DDMMYYYY_ds = "02-01-2006"
	DDMMYYYY_dt = "02.01.2006."

	YYYYMMDD_HHMMSS_sl = "2006/01/02 15:04:05"
	YYYYMMDD_HHMMSS_ds = "2006-01-02 15:04:05"
	YYYYMMDD_HHMMSS_dt = "2006.01.02. 15:04:05"

	DDMMYYYY_HHMMSS_sl = "02/01/2006 15:04:05"
	DDMMYYYY_HHMMSS_ds = "02-01-2006 15:04:05"
	DDMMYYYY_HHMMSS_dt = "02.01.2006. 15:04:05"
)

const DaySeconds = 24 * 60 * 60

var (
	regularYear = [12]int64{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	leapYear    = [12]int64{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

func Systime() int64 {
	return time.Now().Unix()
}

func DateToEpoch(date, format string) int64 {
	t, err := time.Parse(format, date)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return t.Unix()
}

func EpochToDate(e int64, format string) string {
	return time.Unix(e, 0).Format(format)
}

func EpochToDayMonthYear(timestamp int64) (d, m, y int64) {
	datestring := EpochToDate(timestamp, DDMMYYYY_sl)
	parts := strings.Split(datestring, "/")
	d = StringToInt64(parts[0])
	m = StringToInt64(parts[1])
	y = StringToInt64(parts[2])
	return d, m, y
}

func DaysInMonth(year, month int64) int64 {
	if month < 1 || month > 12 {
		return 0
	}
	if IsLeapYear(year) {
		return leapYear[month-1]
	}
	return regularYear[month-1]
}

func IsLeapYear(year int64) bool {
	return year%4 == 0 && !((year%100 == 0) && (year%400 != 0))
}

// FirstDayOfNextMonthEpoch ...
func NextMonths1st(e int64) int64 {
	d, m, y := EpochToDayMonthYear(e)
	m++
	if m > 12 {
		m = 1
		y++
	}
	d = 1

	date := fmt.Sprintf("%02d/%02d/%d", d, m, y)

	return DateToEpoch(date, DDMMYYYY_sl)
}

func SameDate(e1, e2 int64) bool {
	d1, m1, y1 := EpochToDayMonthYear(e1)
	d2, m2, y2 := EpochToDayMonthYear(e2)

	if d1 == d2 && m1 == m2 && y1 == y2 {
		return true
	}

	return false
}

func ToStartOfDay(d int64) int64 {
	rem := d % DaySeconds
	if rem != 0 {
		d -= rem
	}
	return d
}

func ParseTime(date, format string) time.Time {
	t, err := time.Parse(format, date)
	if err != nil {
		fmt.Println(err.Error())
	}
	return t
}
