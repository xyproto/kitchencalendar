//go:build !nb_NO && !en_US

package main

import (
	"errors"
	"time"

	"github.com/xyproto/kal"
)

const msg = "select a locale when building, for example: go build -tags nb_NO"

func formatDate(cal kal.Calendar, date time.Time) string { return msg }
func weekString(week int) string                         { return msg }
func dayAndDate(cal kal.Calendar, t time.Time) string    { return msg }
func newCalendar() (kal.Calendar, error)                 { return nil, errors.New(msg) }
