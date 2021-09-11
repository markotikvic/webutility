package dateutility

import (
	"time"
)

// UnixToDate converts given Unix time to local time in format and returns result:
// YYYY-MM-DD hh:mm:ss +zzzz UTC
func UnixToDate(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// DateToUnix converts given date in Unix timestamp.
func DateToUnix(date interface{}) int64 {
	if date != nil {
		t, ok := date.(time.Time)
		if !ok {
			return 0
		}
		return t.Unix()

	}
	return 0
}

// UnixPtrToDate converts given Unix time to local time in format and returns result:
// YYYY-MM-DD hh:mm:ss +zzzz UTC
func UnixPtrToDatePtr(unix *int64) *time.Time {
	var t time.Time
	if unix == nil {
		return nil
	}
	t = time.Unix(*unix, 0)
	return &t
}

// DateToUnix converts given date in Unix timestamp.
func DatePtrToUnixPtr(date interface{}) *int64 {
	var unix int64

	if date != nil {
		t, ok := date.(time.Time)
		if !ok {
			return nil
		}
		unix = t.Unix()
		return &unix

	}
	return nil
}
