// XXgo:build nb_NO

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/xyproto/kal"
)

// getMonthName takes a time.Time and returns the name of the month in Norwegian
func getMonthName(cal kal.Calendar, t time.Time) string {
	return capitalize(cal.MonthName(t.Month()))
}

// formatDate takes a time.Time and returns a string on the format "17. okt"
func formatDate(cal kal.Calendar, date time.Time) string {
	// Get the day of the month
	day := date.Day()
	// Get the month of the year
	month := date.Month()
	// Get the Norwegian abbreviation for the month
	monthAbbrev := getMonthAbbrev(cal, month)
	// Return the formatted date
	return fmt.Sprintf("%d. %s", day, monthAbbrev)
}

// getMonthAbbrev takes a time.Month and returns the Norwegian abbreviation
func getMonthAbbrev(cal kal.Calendar, month time.Month) string {
	return strings.ToLower(cal.MonthName(month)[:3])
}

// generateWeekHeaderLeft creates the header for the left side of the week table
// on the format "Uke N"
func generateWeekHeaderLeft(year, week int) string {
	return fmt.Sprintf("Uke %d", week)
}

// dayAndDate takes a time.Time and returns the day and date as a string in the form "Mandag 24.".
func dayAndDate(cal kal.Calendar, t time.Time) string {
	// Get the day of the week
	day := t.Weekday()
	// Get the name of the day
	dayName := capitalize(cal.DayName(day))
	// Get the day of the month
	date := t.Day()
	// Return the day and date as a string
	return fmt.Sprintf("%s %d.", dayName, date)
}

// newCalendar returns a new struct that satisfies the kal.Calendar interface
func newCalendar() (kal.Calendar, error) {
	calendar, err := kal.NewCalendar("nb_NO", true)
	if err != nil {
		return nil, err
	}
	return calendar, nil
}
