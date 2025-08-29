package time

import (
	"time"

	timenow "github.com/jinzhu/now"
	"github.com/pkg/errors"
)

const (
	// NowAlias is the string to indicate the current time. It is substituted
	// with 0s when converted to a time object.
	NowAlias string = "now"
)

// ParseString converts a datestamp / duration to a time object.
// @todo add timezone support.
// @todo add support for day interval (not supported by time package).
func ParseString(value string) (time.Time, error) {
	n := time.Now()
	return parseStringWithTime(value, n)
}

// Helper function with injected "now" time object for testing durations.
func parseStringWithTime(value string, now time.Time) (time.Time, error) {
	// Convert "now" string to duration of 0s.
	if value == NowAlias {
		value = "0s"
	}

	// Add our custom formats.
	// k8s datestamp format used in log output.
	timenow.TimeFormats = append(timenow.TimeFormats, time.RFC3339Nano)

	t, err := timenow.New(now).Parse(value)
	if err != nil {
		// Its not an absolute date, try duration (add a "-" to indicate all values being in the past).
		d, err3 := time.ParseDuration("-" + value)
		if err3 != nil {
			return time.Time{}, errors.Wrap(err3, "Could not detect timestamp or duration in string")
		}

		t2 := now.Add(d)
		return t2, nil
	}

	return t, err
}

// UnixMilli returns a Unix timestamp in milliseconds from "January 1, 1970 UTC".
// The result is undefined if the Unix time cannot be represented by an int64.
// Which includes calling UnixMilli on a zero Time is undefined.
//
// This utility is useful for service API's such as CloudWatch Logs which require
// their unix time values to be in milliseconds.
//
// See Go stdlib https://golang.org/pkg/time/#Time.UnixNano for more information.
func UnixMilli(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}
