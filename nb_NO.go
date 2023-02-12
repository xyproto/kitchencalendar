package main

import (
	"fmt"
	"time"

	"github.com/xyproto/kal"
)

// getMonthName takes a time.Time and returns the name of the month in Norwegian
func getMonthName(t time.Time) string {
	// Create a map of month numbers to month names in Norwegian
	monthNames := map[int]string{
		1:  "Januar",
		2:  "Februar",
		3:  "Mars",
		4:  "April",
		5:  "Mai",
		6:  "Juni",
		7:  "Juli",
		8:  "August",
		9:  "September",
		10: "Oktober",
		11: "November",
		12: "Desember",
	}
	// Get the month number from the time.Time
	monthNumber := t.Month()
	// Return the month name from the map
	return monthNames[int(monthNumber)]
}

// formatDate takes a time.Time and returns a string on the format "17. okt"
func formatDate(date time.Time) string {
	// Get the day of the month
	day := date.Day()
	// Get the month of the year
	month := date.Month()
	// Get the Norwegian abbreviation for the month
	monthAbbrev := getMonthAbbrev(month)
	// Return the formatted date
	return fmt.Sprintf("%d. %s", day, monthAbbrev)
}

// getMonthAbbrev takes a time.Month and returns the Norwegian abbreviation
func getMonthAbbrev(month time.Month) string {
	switch month {
	case time.January:
		return "jan"
	case time.February:
		return "feb"
	case time.March:
		return "mar"
	case time.April:
		return "apr"
	case time.May:
		return "mai"
	case time.June:
		return "jun"
	case time.July:
		return "jul"
	case time.August:
		return "aug"
	case time.September:
		return "sep"
	case time.October:
		return "okt"
	case time.November:
		return "nov"
	case time.December:
		return "des"
	default:
		return ""
	}
}

// generateWeekHeaderLeft creates the header for the left side of the week table
// on the format "Uke N"
func generateWeekHeaderLeft(year, week int) string {
	return fmt.Sprintf("Uke %d", week)
}

// dayAndDate takes a time.Time and returns the day and date as a string in the form "Mandag 24.".
func dayAndDate(t time.Time) string {
	// Get the day of the week
	day := t.Weekday().String()
	// Get the day of the month
	date := t.Day()
	// Map the day of the week to Norwegian
	dayMap := map[string]string{
		"Monday":    "Mandag",
		"Tuesday":   "Tirsdag",
		"Wednesday": "Onsdag",
		"Thursday":  "Torsdag",
		"Friday":    "Fredag",
		"Saturday":  "Lørdag",
		"Sunday":    "Søndag",
	}
	// Return the day and date as a string
	return fmt.Sprintf("%s %d.", dayMap[day], date)
}

func NewCalendar() (*kal.Calendar, error) {
	calendar, err := kal.NewCalendar("nb_NO", true)
	if err != nil {
		return nil, err
	}
	return &calendar, nil
}
