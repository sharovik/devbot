package time

import (
	"time"

	"github.com/sharovik/devbot/internal/config"
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
