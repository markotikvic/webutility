package webutility

import (
	"fmt"
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

// EqualQuotes encapsulates given string in SQL 'equal' statement and returns result.
// Example: "hello" -> " = 'hello'"
func EqualQuotes(stmt string) string {
	if stmt != "" {
		stmt = fmt.Sprintf(" = '%s'", stmt)
	}
	return stmt
}

// EqualString ...
func EqualString(stmt string) string {
	if stmt != "" {
		stmt = fmt.Sprintf(" = %s", stmt)
	}
	return stmt
}

// LikeQuotes encapsulates given string in SQL 'like' statement and returns result.
// Example: "hello" -> " LIKE UPPER('%hello%')"
func LikeQuotes(stmt string) string {
	if stmt != "" {
		stmt = fmt.Sprintf(" LIKE UPPER('%s%s%s')", "%", stmt, "%")
	}
	return stmt
}
