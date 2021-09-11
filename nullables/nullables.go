package nullables

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	date "markotikvic/webutility/dateutility"
)

// NullBool is a wrapper for sql.NullBool with added JSON (un)marshalling
type NullBool sql.NullBool

// Scan ...
func (nb *NullBool) Scan(value interface{}) error {
	var b sql.NullBool
	if err := b.Scan(value); err != nil {
		nb.Bool, nb.Valid = false, false
		return err
	}
	nb.Bool, nb.Valid = b.Bool, b.Valid
	return nil
}

// Value ...
func (nb *NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

// MarshalJSON ...
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON ...
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	var temp *bool
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	if temp != nil {
		nb.Valid = true
		nb.Bool = *temp
	} else {
		nb.Valid = false
	}
	return nil
}

// CastToSQL ...
func (nb *NullBool) CastToSQL() sql.NullBool {
	return sql.NullBool(*nb)
}

// NullString is a wrapper for sql.NullString with added JSON (un)marshalling
type NullString sql.NullString

// Scan ...
func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		ns.String, ns.Valid = "", false
		return err
	}
	ns.String, ns.Valid = s.String, s.Valid
	return nil
}

// Value ...
func (ns *NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// Val ...
func (ns *NullString) Val() string {
	return ns.String
}

// MarshalJSON ...
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON ...
func (ns *NullString) UnmarshalJSON(b []byte) error {
	var temp *string
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	if temp != nil {
		ns.Valid = true
		ns.String = *temp
	} else {
		ns.Valid = false
	}
	return nil
}

// CastToSQL ...
func (ns *NullString) CastToSQL() sql.NullString {
	return sql.NullString(*ns)
}

// NullInt64 is a wrapper for sql.NullInt64 with added JSON (un)marshalling
type NullInt64 sql.NullInt64

// Scan ...
func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		ni.Int64, ni.Valid = 0, false
		return err
	}
	ni.Int64, ni.Valid = i.Int64, i.Valid
	return nil
}

// ScanPtr ...
func (ni *NullInt64) ScanPtr(v interface{}) error {
	if ip, ok := v.(*int64); ok && ip != nil {
		return ni.Scan(*ip)
	}
	return nil
}

// Value ...
func (ni *NullInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}

func (ni *NullInt64) Val() int64 {
	return ni.Int64
}

// Add
func (ni *NullInt64) Add(i NullInt64) {
	ni.Valid = true
	ni.Int64 += i.Int64
}

func (ni *NullInt64) Set(i int64) {
	ni.Valid = true
	ni.Int64 = i
}

// MarshalJSON ...
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON ...
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	var temp *int64
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	if temp != nil {
		ni.Valid = true
		ni.Int64 = *temp
	} else {
		ni.Valid = false
	}
	return nil
}

// CastToSQL ...
func (ni *NullInt64) CastToSQL() sql.NullInt64 {
	return sql.NullInt64(*ni)
}

// NullFloat64 is a wrapper for sql.NullFloat64 with added JSON (un)marshalling
type NullFloat64 sql.NullFloat64

// Scan ...
func (nf *NullFloat64) Scan(value interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(value); err != nil {
		nf.Float64, nf.Valid = 0.0, false
		return err
	}
	nf.Float64, nf.Valid = f.Float64, f.Valid
	return nil
}

// ScanPtr ...
func (nf *NullFloat64) ScanPtr(v interface{}) error {
	if fp, ok := v.(*float64); ok && fp != nil {
		return nf.Scan(*fp)
	}
	return nil
}

// Value ...
func (nf *NullFloat64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}

// Val ...
func (nf *NullFloat64) Val() float64 {
	return nf.Float64
}

// Add ...
func (nf *NullFloat64) Add(f NullFloat64) {
	nf.Valid = true
	nf.Float64 += f.Float64
}

func (nf *NullFloat64) Set(f float64) {
	nf.Valid = true
	nf.Float64 = f
}

// MarshalJSON ...
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON ...
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	var temp *float64
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	if temp != nil {
		nf.Valid = true
		nf.Float64 = *temp
	} else {
		nf.Valid = false
	}
	return nil
}

// CastToSQL ...
func (nf *NullFloat64) CastToSQL() sql.NullFloat64 {
	return sql.NullFloat64(*nf)
}

// NullDateTime ...
type NullDateTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan ...
func (nt *NullDateTime) Scan(value interface{}) (err error) {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		return
	case []byte:
		nt.Time, err = parseSQLDateTime(string(v), time.UTC)
		nt.Valid = (err == nil)
		return
	case string:
		nt.Time, err = parseSQLDateTime(v, time.UTC)
		nt.Valid = (err == nil)
		return
	}

	nt.Valid = false
	return fmt.Errorf("Can't convert %T to time.Time", value)
}

// Value implements the driver Valuer interface.
func (nt NullDateTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON ...
func (nt NullDateTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		format := nt.Time.Format("2006-01-02 15:04:05")
		return json.Marshal(format)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON ...
func (nt *NullDateTime) UnmarshalJSON(b []byte) error {
	var temp *time.Time
	var t1 time.Time
	var err error

	s1 := string(b)
	s2 := s1[1 : len(s1)-1]
	if s1 == "null" {
		temp = nil
	} else {
		t1, err = time.Parse("2006-01-02 15:04:05", s2)
		if err != nil {
			return err
		}
		temp = &t1
	}

	if temp != nil {
		nt.Valid = true
		nt.Time = *temp
	} else {
		nt.Valid = false
	}
	return nil
}

func (nt *NullDateTime) CastToSQL() NullDateTime {
	return *nt
}

func parseSQLDateTime(str string, loc *time.Location) (t time.Time, err error) {
	base := "0000-00-00 00:00:00.0000000"
	timeFormat := "2006-01-02 15:04:05.999999"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("invalid time string: %s", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return
}

// NullDate ...
type NullDate struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan ...
func (nt *NullDate) Scan(value interface{}) (err error) {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		return
	case []byte:
		nt.Time, err = parseSQLDate(string(v), time.UTC)
		nt.Valid = (err == nil)
		return
	case string:
		nt.Time, err = parseSQLDate(v, time.UTC)
		nt.Valid = (err == nil)
		return
	}

	nt.Valid = false
	return fmt.Errorf("Can't convert %T to time.Time", value)
}

// Value implements the driver Valuer interface.
func (nt NullDate) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON ...
func (nt NullDate) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		format := nt.Time.Format("2006-01-02")
		return json.Marshal(format)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON ...
func (nt *NullDate) UnmarshalJSON(b []byte) error {
	var temp *time.Time
	var t1 time.Time
	var err error

	s1 := string(b)
	s2 := s1[1 : len(s1)-1]
	if s1 == "null" {
		temp = nil
	} else {
		t1, err = time.Parse("2006-01-02", s2)
		if err != nil {
			return err
		}
		temp = &t1
	}

	if temp != nil {
		nt.Scan(t1)
	} else {
		nt.Valid = false
	}
	return nil
}

func (nt *NullDate) CastToSQL() NullDate {
	return *nt
}

func (nd *NullDate) Format(f string) string {
	if !nd.Valid {
		return ""
	}

	return date.EpochToDate(nd.Time.Unix(), f)
}

func parseSQLDate(str string, loc *time.Location) (t time.Time, err error) {
	base := "0000-00-00"
	timeFormat := "2006-01-02"
	switch len(str) {
	case 10:
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("invalid time string: %s", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return
}

// NullTime ...
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan ...
func (nt *NullTime) Scan(value interface{}) (err error) {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		return
	case []byte:
		nt.Time, err = parseSQLTime(string(v), time.UTC)
		nt.Valid = (err == nil)
		return
	case string:
		nt.Time, err = parseSQLTime(v, time.UTC)
		nt.Valid = (err == nil)
		return
	}

	nt.Valid = false
	return fmt.Errorf("Can't convert %T to time.Time", value)
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON ...
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		format := nt.Time.Format("15:04:05")
		return json.Marshal(format)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON ...
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	var temp *time.Time
	var t1 time.Time
	var err error

	s1 := string(b)
	s2 := s1[1 : len(s1)-1]
	if s1 == "null" {
		temp = nil
	} else {
		t1, err = time.Parse("2006-05-04 15:04:05", "1970-01-01 "+s2)
		if err != nil {
			return err
		}
		temp = &t1
	}

	if temp != nil {
		nt.Scan(t1)
	} else {
		nt.Valid = false
	}
	return nil
}

func (nt *NullTime) CastToSQL() NullTime {
	return *nt
}

// NOTE(marko): Date must be included because database can't convert it to TIME otherwise.
func parseSQLTime(str string, loc *time.Location) (t time.Time, err error) {
	base := "00:00:00"
	timeFormat := "15:04:05"
	switch len(str) {
	case 8:
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse("2006-05-04 "+timeFormat[:len(str)], "1970-01-01 "+str)
	default:
		err = fmt.Errorf("invalid time string: %s", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return
}
