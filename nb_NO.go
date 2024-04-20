//go:build nb_NO

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
	monthAbbrev := GetMonthAbbrev(cal, month)
	// Return the formatted date
	return fmt.Sprintf("%d. %s", day, monthAbbrev)
}

// WeekString creates the header for the left side of the week table
func WeekString(week int) string {
	return fmt.Sprintf("Uke %d", week)
}

// DayAndDate takes a time.Time and returns the day and date as a string in the form "Mandag 24.".
func DayAndDate(cal kal.Calendar, t time.Time) string {
	// Get the day of the week
	day := t.Weekday()
	// Get the name of the day
	dayName := capitalize(cal.DayName(day))
	// Get the day of the month
	date := t.Day()
	// Return the day and date as a string
	return fmt.Sprintf("%s %d.", dayName, date)
}

// NewCalendar returns a new struct that satisfies the kal.Calendar interface
func NewCalendar() (kal.Calendar, error) {
	calendar, err := kal.NewCalendar("nb_NO", true)
	if err != nil {
		return nil, err
	}
	return calendar, nil
}
