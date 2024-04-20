//go:build !nb_NO && !en_US

package kitchencalendar

import (
	"errors"
	"time"

	"github.com/xyproto/kal"
)

const msg = "select a locale when building, for example: go build -tags nb_NO"

func FormatDate(cal kal.Calendar, date time.Time) string { return msg }
func WeekString(week int) string                         { return msg }
func DayAndDate(cal kal.Calendar, t time.Time) string    { return msg }
func NewCalendar() (kal.Calendar, error)                 { return nil, errors.New(msg) }
