//go:build en_US

package kitchencalendar

import (
	"fmt"
	"time"

	"github.com/xyproto/kal"
)

// FormatDate takes a time.Time and returns a string on the format "17. okt"
func FormatDate(cal kal.Calendar, date time.Time) string {
	// Get the day of the month
	day := date.Day()
	// Get the month of the year
	month := date.Month()
	// Get the calendar abbreviation for the month
	monthAbbrev := capitalize(GetMonthAbbrev(cal, month))
	// Create the suffix for the date
	suffix := "th"
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	// Return the formatted date
	return fmt.Sprintf("%d%s of %s", day, suffix, monthAbbrev)
}

// WeekString creates the header for the left side of the week table
func WeekString(week int) string {
	return fmt.Sprintf("Week %d", week)
}

// DayAndDate takes a time.Time and returns the day and date as a string
// on the form "Mon. 24st"
func DayAndDate(cal kal.Calendar, t time.Time) string {
	// Get the day of the week
	dayName := t.Weekday().String()
	// Abbreviate the day
	dayName = dayName[:3]
	// Get the day of the month
	day := t.Day()
	// Create the suffix for the date
	suffix := "th"
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	// Return the day and date string
	return fmt.Sprintf("%s. %d%s", dayName, day, suffix)
}

// NewCalendar returns a new struct that satisfies the kal.Calendar interface
func NewCalendar() (kal.Calendar, error) {
	calendar, err := kal.NewCalendar("en_US", true)
	if err != nil {
		return nil, err
	}
	return calendar, nil
}
