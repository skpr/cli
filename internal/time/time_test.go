package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestValue struct {
	Value    string
	Expected string
	Comment  string
}

func TestParseStringWithTimezone(t *testing.T) {
	TestValues := []TestValue{
		{
			Value:    "2017-6-15",
			Expected: "2017-06-15T00:00:00Z",
			Comment:  "date only",
		},
		{
			Value:    "19:05",
			Expected: "2019-07-01T19:05:00Z",
			Comment:  "time only",
		},
		{
			Value:    "2017-6-15 19:05",
			Expected: "2017-06-15T19:05:00Z",
			Comment:  "date and time",
		},
		{
			Value:    "30s",
			Expected: "2019-07-01T08:59:30Z",
			Comment:  "relative seconds",
		},
		{
			Value:    "15m",
			Expected: "2019-07-01T08:45:00Z",
			Comment:  "relative minutes",
		},
		{
			Value:    "2h",
			Expected: "2019-07-01T07:00:00Z",
			Comment:  "relative hours",
		},
		{
			Value:    "1w",
			Expected: "2019-06-24T09:00:00Z",
			Comment:  "relative days",
		},
		{
			Value:    "2d4h",
			Expected: "2019-06-29T05:00:00Z",
			Comment:  "relative days",
		},
		{
			Value:    "now",
			Expected: "2019-07-01T09:00:00Z",
			Comment:  "now",
		},
		{
			Value:    "2019-04-26T07:46:40.795968088Z",
			Expected: "2019-04-26T07:46:40Z",
			Comment:  "kubernetes datestamp",
		},
	}

	// Use a static date for "now" so durations / relative dates can be consistently tested.
	now, err := time.Parse(time.RFC3339, "2019-07-01T09:00:00Z")
	assert.NoError(t, err)
	for _, item := range TestValues {
		actual, err := parseStringWithTime(item.Value, now)
		assert.NoError(t, err)
		assert.Equal(t, item.Expected, actual.Format(time.RFC3339), item.Comment)
	}
}
