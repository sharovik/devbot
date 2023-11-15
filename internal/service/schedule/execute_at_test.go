package schedule

import (
	"testing"
	"time"

	_time "github.com/sharovik/devbot/internal/service/time"

	"github.com/stretchr/testify/assert"
)

func TestExecuteAt_IsEmpty(t *testing.T) {
	cases := []ExecuteAt{
		{
			Minutes: 1,
		},
		{
			Hours: 1,
		},
		{
			ExactDatetime: _time.Service.Now(),
		},
	}

	for _, actual := range cases {
		assert.False(t, actual.IsEmpty())
	}

	assert.True(t, new(ExecuteAt).IsEmpty())
}

func TestExecuteAt_parseDays(t *testing.T) {
	var (
		cases = map[string]int64{
			"in 2 days":   int64(2),
			"after 1 day": int64(1),
			"every day":   int64(0),
		}
	)

	for text, expected := range cases {
		e := ExecuteAt{}
		err := e.parseDays(text)
		assert.NoError(t, err)
		assert.Equal(t, expected, e.Days)
	}
}

func TestExecuteAt_parse(t *testing.T) {
	var (
		cases = map[string]interface{}{
			"in 2 hours":   2,
			"after 1 hour": 1,
			"every hour":   nil,
		}
	)

	for text, expected := range cases {
		e := ExecuteAt{}
		hours, err := e.parse(text, HourRegexp)
		assert.NoError(t, err)
		assert.Equal(t, expected, hours)
	}

	cases = map[string]interface{}{
		"in 20 minutes": 20,
		"every minute":  nil,
	}

	for text, expected := range cases {
		e := ExecuteAt{}
		minutes, err := e.parse(text, MinuteRegexp)
		assert.NoError(t, err)
		assert.Equal(t, expected, minutes)
	}
}

func TestExecuteAt_isRepeatable(t *testing.T) {
	var (
		cases = map[string]bool{
			"repeat daily": true,
			"every hour":   true,
			"repeat 1 days and 9 hours and 30 minutes": true,
			"in 1 hour": false,
		}
	)

	for text, expected := range cases {
		e := ExecuteAt{}
		assert.Equal(t, expected, e.isRepeatable(text), text)
	}
}

func TestExecuteAt_isDelayed(t *testing.T) {
	var (
		cases = map[string]bool{
			"in few hours": true,
			"every hour":   false,
			"after 1 hour": true,
		}
	)

	for text, expected := range cases {
		e := ExecuteAt{}
		assert.Equal(t, expected, e.isDelayed(text), text)
	}
}

func TestExecuteAt_getDatetime(t *testing.T) {
	var (
		actual       ExecuteAt
		expectedDate time.Time
		ct           = _time.Service.Now()
		err          error
	)

	actual, err = new(ExecuteAt).FromString("schedule event examplescenario every 1 minute")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute()+1, 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))

	actual, err = new(ExecuteAt).FromString("in 1 hour and 2 minutes")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), ct.Day(), ct.Hour()+1, ct.Minute()+2, 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))

	actual, err = new(ExecuteAt).FromString("1 hour")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute()+60, 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))

	actual, err = new(ExecuteAt).FromString("23 minutes")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute()+23, 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))

	actual, err = new(ExecuteAt).FromString("2022-12-18 11:22")
	assert.NoError(t, err)
	expectedDate, err = time.Parse(timeFormat, "2022-12-18 11:22")
	assert.NoError(t, err)
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))

	actual, err = new(ExecuteAt).FromString("in 20 minutes")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute()+20, 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))

	actual, err = new(ExecuteAt).FromString("in 1 day")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), ct.Day(), ct.Hour()+24, ct.Minute(), 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))

	actual, err = new(ExecuteAt).FromString("2 days")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), ct.Day()+2, ct.Hour(), ct.Minute(), 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))

	actual, err = new(ExecuteAt).FromString("repeat 1 days and at 9:30")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), ct.Day()+1, 9, 30, 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))
	assert.Equal(t, "repeat 1 days and at 9:30", actual.toString())
}

func TestExecuteAt_ParseDays(t *testing.T) {
	var (
		actual       ExecuteAt
		expectedDate time.Time
		ct           = _time.Service.Now()
		err          error
	)

	days := int((7 + (time.Sunday - ct.Weekday())) % 7)
	now := _time.Service.Now()
	_, _, d := now.AddDate(0, 0, days).Date()

	actual, err = new(ExecuteAt).FromString("Sunday at 10:00")
	assert.NoError(t, err)
	expectedDate = time.Date(ct.Year(), ct.Month(), d, 10, 00, 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))
	assert.Equal(t, "Sunday and at 10:0", actual.toString())

	actual, err = new(ExecuteAt).FromString("every monday at 9:10")
	assert.NoError(t, err)

	days = int((7 + (time.Monday - ct.Weekday())) % 7)
	now = _time.Service.Now()
	_, _, d = now.AddDate(0, 0, days).Date()
	expectedDate = time.Date(ct.Year(), ct.Month(), d, 9, 10, 0, 0, ct.Location())
	assert.Equal(t, expectedDate.Format(timeFormat), actual.getDatetime().Format(timeFormat))
	assert.Equal(t, "repeat Monday and at 9:10", actual.toString())
}
