package time

import (
	"github.com/sharovik/devbot/internal/config"
	"time"
)

var Service TimeService

type TimeService struct {
	TimeZone *time.Location
}

func (s TimeService) Now() time.Time {
	if s.TimeZone == nil {
		s.TimeZone = config.DefaultTimezone
	}

	return time.Now().In(s.TimeZone)
}

func InitNOW(timeZone *time.Location) {
	Service = TimeService{
		TimeZone: timeZone,
	}
}
