package utils

import (
	"fmt"
	"strconv"
	"time"
)

var strToTimeFormats = []string{
	"2006-01-02 15:04:05 Z0700 MST",
	"2006-01-02 15:04:05 Z07:00 MST",
	"2006-01-02 15:04:05 Z0700 -0700",
	"Mon Jan _2 15:04:05 -0700 MST 2006",
	time.RFC822Z, // "02 Jan 06 15:04 -0700"
	"01/02/06",
	"01/02/2006",
	"01/02/2006 15:04:05",
	"01/02/2006 15:04",
	"2006/01/02",
	"2006-01-02 15:04",
	"2006-01-02T15:04",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02",
	time.RFC3339, // "2006-01-02T15:04:05Z07:00", RFC3339Nano"
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05 Z07:00",
	time.RubyDate, // "Mon Jan 02 15:04:05 -0700 2006"
	time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
	time.RFC822,   // "02 Jan 06 15:04 MST"
	"2006-01-02 15:04:05 MST",
	time.UnixDate, // "Mon Jan _2 15:04:05 MST 2006"
	time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC850,   // "Monday, 02-Jan-06 15:04:05 MST"
	time.Kitchen,  // "3:04PM"
	time.Stamp,    // "Jan _2 15:04:05", time.StampMilli, time.StampMicro, time.StampNano"
	time.ANSIC,    // "Mon Jan _2 15:04:05 2006"
	"Jan _2, 2006",
	"01/02/06 15:04",
	"01/02/06 15:04:05",
	"_2/Jan/2006 15:04:05",
}

// ToTime try to convert the argument into a time
func ToTime(val interface{}) time.Time {
	str := ToString(val)
	for _, layout := range strToTimeFormats {
		r, err := time.ParseInLocation(layout, str, time.Local)
		if err == nil {
			return r
		}
	}
	return time.Now()
}

// ToString try to convert the argument into a string
func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

// ToInt try to convert the argument into a int
func ToInt(val interface{}) int {
	str := ToString(val)
	r, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		r = 0
	}
	return int(r)
}

// ToInt64 try to convert the argument into a int64
func ToInt64(val interface{}) int64 {
	str := ToString(val)
	r, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		r = 0
	}
	return r
}

// ToUint32 try to convert the argument into a uint32
func ToUint32(val interface{}) uint32 {
	str := ToString(val)
	r, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		r = 0
	}
	return uint32(r)
}

// ToUint64 try to convert the argument into a uint64
func ToUint64(val interface{}) uint64 {
	str := ToString(val)
	r, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		r = 0
	}
	return r
}

// ToFloat64 try to convert the argument into a float64
func ToFloat64(val interface{}) float64 {
	str := ToString(val)
	r, err := strconv.ParseFloat(str, 64)
	if err != nil {
		r = 0
	}
	return r
}

// ToBool try to convert the argument into a bool
func ToBool(val interface{}) bool {
	r, err := strconv.ParseBool(ToString(val))
	if err != nil {
		r = false
	}
	return r
}
