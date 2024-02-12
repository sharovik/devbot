package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	_time "github.com/sharovik/devbot/internal/service/time"

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
	exactTimeRegexp   = `(?im)(\d+):(\d+)`

	timeFormat = "2006-01-02 15:04"
)

var daysOfWeek = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

type ExecuteAt struct {
	Days          int64
	Minutes       int64
	Hours         int64
	Weekday       interface{}
	IsRepeatable  bool
	IsDelayed     bool
	ExactDatetime time.Time
	IsExactHours  bool
}

func (e *ExecuteAt) parseExactTime(text string) error {
	res := helper.FindMatches(exactTimeRegexp, text)
	if len(res) == 0 {
		return nil
	}

	hour, err := strconv.Atoi(res["1"])
	if err != nil {
		return err
	}

	minute, err := strconv.Atoi(res["2"])
	if err != nil {
		return err
	}

	e.Hours = int64(hour)
	e.Minutes = int64(minute)

	e.IsExactHours = true

	return nil
}

func (e *ExecuteAt) getDatetime() time.Time {
	t := _time.Service.Now()

	if e.Days != 0 || e.Minutes != 0 || e.Hours != 0 {
		hours := t.Hour()
		if e.Hours != 0 {
			hours = int(e.Hours)
		}

		minutes := int(e.Minutes)
		if !e.IsExactHours {
			minutes = t.Minute() + minutes
		}

		if e.IsRepeatable || e.IsDelayed {
			e.generateDelayedDate()
			return e.ExactDatetime
		}

		return time.Date(t.Year(), t.Month(), e.generateDays(t), hours, minutes, 0, 0, t.Location())
	}

	return e.ExactDatetime
}

func (e *ExecuteAt) generateDays(now time.Time) int {
	if e.Weekday == nil {
		days := now.Day()
		if e.Days != 0 {
			days += int(e.Days)
		}

		return days
	}

	if e.Weekday.(time.Weekday) == now.Weekday() {
		return now.Day()
	}

	days := int((7 + (e.Weekday.(time.Weekday) - now.Weekday())) % 7)
	_, _, d := now.AddDate(0, 0, days).Date()
	return d
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

	if e.Weekday != nil {
		res = append(res, e.Weekday.(time.Weekday).String())
	}

	if e.IsExactHours {
		res = append(res, fmt.Sprintf("at %d:%d", e.Hours, e.Minutes))
	} else {
		if e.Hours != 0 {
			res = append(res, fmt.Sprintf("%d hours", e.Hours))
		}

		if e.Minutes != 0 {
			res = append(res, fmt.Sprintf("%d minutes", e.Minutes))
		}
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

	result, err := time.ParseInLocation(timeFormat, text, _time.Service.TimeZone)
	if err != nil {
		return err
	}

	e.ExactDatetime = result
	return nil
}

func (e *ExecuteAt) parse(text string, regex string) (result interface{}, err error) {
	res := helper.FindMatches(regex, text)
	if res["2"] == "" {
		return nil, nil
	}

	return strconv.Atoi(res["2"])
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

func (e *ExecuteAt) parseWeekday(text string) error {
	var days []string
	e.Weekday = nil
	for dayName := range daysOfWeek {
		days = append(days, dayName)
	}

	regexStr := fmt.Sprintf("(?i)(%s)", strings.Join(days, "|"))
	res := helper.FindMatches(regexStr, text)
	if res["1"] == "" {
		return nil
	}

	dayName := strings.ToLower(res["1"])

	e.Weekday = daysOfWeek[dayName]

	return nil
}

func (e *ExecuteAt) parseHoursAndMinutes(text string) error {
	var (
		hours   interface{}
		minutes interface{}
		err     error
	)

	if hours, err = e.parse(text, HourRegexp); err != nil {
		return err
	}

	if minutes, err = e.parse(text, MinuteRegexp); err != nil {
		return err
	}

	if hours == nil && minutes == nil {
		return nil
	}

	//When we receive only hours but not minutes, we convert hours as minutes
	if hours != nil && minutes == nil {
		e.Minutes = int64(hours.(int) * 60)

		return nil
	}

	if hours == nil {
		e.Hours = 0
	} else {
		e.Hours = int64(hours.(int))
	}

	e.Minutes = int64(minutes.(int))

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

	if err := e.parseHoursAndMinutes(text); err != nil {
		return ExecuteAt{}, err
	}

	if err := e.parseWeekday(text); err != nil {
		return ExecuteAt{}, err
	}

	if err := e.parseDays(text); err != nil {
		return ExecuteAt{}, err
	}

	if err := e.parseExactTime(text); err != nil {
		return ExecuteAt{}, err
	}

	if e.IsDelayed || e.IsRepeatable {
		e.generateDelayedDate()
	}

	return *e, nil
}

func (e *ExecuteAt) generateDelayedDate() {
	t := _time.Service.Now()
	days := t.Day()
	if e.Days != 0 {
		days += int(e.Days)
	}

	hours := t.Hour()
	minutes := t.Minute()

	if !e.IsExactHours {
		if e.Hours != 0 {
			hours += int(e.Hours)
		}

		if e.Minutes != 0 {
			minutes += int(e.Minutes)
		}
	} else {
		hours = int(e.Hours)
		minutes = int(e.Minutes)
	}

	e.ExactDatetime = time.Date(t.Year(), t.Month(), e.generateDays(t), hours, minutes, 0, 0, t.Location())
}
