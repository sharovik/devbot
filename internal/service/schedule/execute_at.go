package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sharovik/devbot/internal/helper"
)

const (
	//DatetimeRegexp regexp for datetime parsing
	DatetimeRegexp = `(?im)(\d+-\d+-\d+ \d+:\d+)`

	//MinuteRegexp regexp for minutes parsing
	MinuteRegexp = `(?im)((\d+) minute|minutes)`

	//HourRegexp regexp for hours parsing
	HourRegexp = `(?im)((\d+) hour|hours)`

	//DayRegexp regexp for hours parsing
	DayRegexp = `(?im)((\d+) day|days)`

	repeatableRegexp  = `(?im)(?:^|\s)(repeat|every)\s`
	delayedTimeRegexp = `(?im)(?:^|\s)(in|after)\s`

	timeFormat = "2006-01-02 15:04"
)

type ExecuteAt struct {
	Days          int64
	Minutes       int64
	Hours         int64
	IsRepeatable  bool
	IsDelayed     bool
	ExactDatetime time.Time
}

func (e *ExecuteAt) getDatetime() time.Time {
	t := time.Now()

	if e.Days != 0 || e.Minutes != 0 || e.Hours != 0 {
		days := t.Day()
		if e.Days != 0 {
			days = int(e.Days)
		}

		hours := t.Hour()
		if e.Hours != 0 {
			hours = int(e.Hours)
		}

		minutes := t.Minute()
		if e.Minutes != 0 {
			minutes = int(e.Minutes)
		}

		if e.IsRepeatable || e.IsDelayed {
			e.generateDelayedDate()
			return e.ExactDatetime
		}

		return time.Date(t.Year(), t.Month(), days, hours, minutes, 0, 0, time.UTC)
	}

	return e.ExactDatetime
}

func (e *ExecuteAt) IsEmpty() bool {
	return e.Days == 0 && e.Hours == 0 && e.Minutes == 0 && e.ExactDatetime.IsZero()
}

func (e *ExecuteAt) toString() string {
	if e.IsEmpty() {
		return ""
	}

	if !e.ExactDatetime.IsZero() && !e.IsRepeatable {
		return e.ExactDatetime.Format(timeFormat)
	}

	var res []string
	if e.Days != 0 {
		res = append(res, fmt.Sprintf("%d days", e.Days))
	}

	if e.Hours != 0 {
		res = append(res, fmt.Sprintf("%d hours", e.Hours))
	}

	if e.Minutes != 0 {
		res = append(res, fmt.Sprintf("%d minutes", e.Minutes))
	}

	if len(res) == 0 {
		return ""
	}

	result := ""
	if e.IsRepeatable {
		result = "repeat "
	}

	return fmt.Sprintf("%s%s", result, strings.Join(res, " and "))
}

func (e *ExecuteAt) parseDateTime(text string) error {
	res := helper.FindMatches(DatetimeRegexp, text)
	if res["1"] == "" {
		return nil
	}

	result, err := time.ParseInLocation(timeFormat, text, time.UTC)
	if err != nil {
		return err
	}

	e.ExactDatetime = result
	return nil
}

func (e *ExecuteAt) parseMinutes(text string) error {
	res := helper.FindMatches(MinuteRegexp, text)
	if res["2"] == "" {
		return nil
	}

	minutes, err := strconv.Atoi(res["2"])
	if err != nil {
		return err
	}

	e.Minutes = int64(minutes)

	return nil
}

func (e *ExecuteAt) parseDays(text string) error {
	res := helper.FindMatches(DayRegexp, text)
	if res["2"] == "" {
		return nil
	}

	days, err := strconv.Atoi(res["2"])
	if err != nil {
		return err
	}

	e.Days = int64(days)

	return nil
}

func (e *ExecuteAt) parseHours(text string) error {
	res := helper.FindMatches(HourRegexp, text)
	if res["2"] == "" {
		return nil
	}

	hours, err := strconv.Atoi(res["2"])
	if err != nil {
		return err
	}

	e.Hours = int64(hours)
	return nil
}

func (e *ExecuteAt) isRepeatable(text string) bool {
	res := helper.FindMatches(repeatableRegexp, text)

	return res["1"] != ""
}

func (e *ExecuteAt) isDelayed(text string) bool {
	res := helper.FindMatches(delayedTimeRegexp, text)

	return res["1"] != ""
}

func (e *ExecuteAt) FromString(text string) (ExecuteAt, error) {
	if err := e.parseDateTime(text); err != nil {
		return ExecuteAt{}, err
	}

	if !e.IsEmpty() {
		return *e, nil
	}

	if e.isRepeatable(text) {
		e.IsRepeatable = true
	}

	if e.isDelayed(text) {
		e.IsDelayed = true
	}

	if err := e.parseHours(text); err != nil {
		return ExecuteAt{}, err
	}

	if err := e.parseDays(text); err != nil {
		return ExecuteAt{}, err
	}

	if err := e.parseMinutes(text); err != nil {
		return ExecuteAt{}, err
	}

	if e.IsDelayed || e.IsRepeatable {
		e.generateDelayedDate()
	}

	return *e, nil
}

func (e *ExecuteAt) generateDelayedDate() {
	t := time.Now()
	days := t.Day()
	if e.Days != 0 {
		days += int(e.Days)
	}

	hours := t.Hour()
	if e.Hours != 0 {
		hours += int(e.Hours)
	}

	minutes := t.Minute()
	if e.Minutes != 0 {
		minutes += int(e.Minutes)
	}

	e.ExactDatetime = time.Date(t.Year(), t.Month(), days, hours, minutes, 0, 0, time.UTC)
}
