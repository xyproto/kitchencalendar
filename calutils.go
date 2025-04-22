package kitchencalendar

import (
	"strings"
	"time"

	"github.com/xyproto/kal"
)

// FirstMondayOfWeek finds the first monday of the week, given a year and a week number
// based on ISO 8601 standard
func FirstMondayOfWeek(year, week int) time.Time {
	// The ISO 8601 definition: week 1 is the week with the first Thursday of the year
	// This means week 1 contains January 4
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)
	// Find the Monday of the week containing Jan 4
	firstWeekMonday := jan4
	for firstWeekMonday.Weekday() != time.Monday {
		firstWeekMonday = firstWeekMonday.AddDate(0, 0, -1)
	}
	// Now we can calculate the Monday of the specified week
	weekOffset := (week - 1) * 7
	return firstWeekMonday.AddDate(0, 0, weekOffset) // monday
}

// FirstSundayOfWeek finds the first Sunday of the week, given a year and a week number
// based on ISO 8601 standard
func FirstSundayOfWeek(year, week int) time.Time {
	// Get the Monday of the specified week
	monday := FirstMondayOfWeek(year, week)
	// Go forward 6 days to get to the Sunday of that week
	return monday.AddDate(0, 0, 6) // sunday
}

// FirstSundayAfter finds the first Sunday after the given date
func FirstSundayAfter(date time.Time) time.Time {
	// Get the day of the week for the given date
	dayOfWeek := date.Weekday()
	// Calculate the number of days to add to the given date to get the first Sunday
	daysToAdd := 7 - int(dayOfWeek)
	// Add the calculated number of days to the given date and return it
	return date.AddDate(0, 0, daysToAdd)
}

// FirstSaturdayAfter takes a time.Time and returns the first Saturday after the given date as a time.Time
func FirstSaturdayAfter(date time.Time) time.Time {
	// Get the day of the week for the given date
	dayOfWeek := date.Weekday()
	// Calculate the number of days until the next Saturday
	daysToAdd := 6 - int(dayOfWeek)
	// Add the number of days until the next Saturday to the given date
	return date.AddDate(0, 0, daysToAdd)
}

// IterateDays iterates over days from startDay to endDay (inclusive) and calls f for each day
func IterateDays(startDay, endDay time.Time, f func(time.Time) error) error {
	// Create a new time.Time object representing the start of the startDay
	start := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	// Create a new time.Time object representing the start of the endDay
	end := time.Date(endDay.Year(), endDay.Month(), endDay.Day(), 0, 0, 0, 0, endDay.Location())
	// Iterate over the days from start to end
	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		// Call the function with the current day
		if err := f(d); err != nil {
			return err
		}
	}
	return nil
}

// GetCurrentYear returns the current year as an int
func GetCurrentYear() int {
	return time.Now().Year()
}

// GetCurrentWeek returns the current week number as an int
func GetCurrentWeek() int {
	// Get the current time
	now := time.Now()
	// Get the ISO year and week number
	currentYear, week := now.ISOWeek()
	// Make sure we're using weeks for the current year, not next year's week 1
	if currentYear != now.Year() {
		// If the ISO year doesn't match the calendar year, we need to adjust
		return GetWeekForDate(now)
	}
	// Return the week number
	return week
}

// GetWeekForDate returns the week number for a given date, ensuring it's
// attributed to the correct year (not the ISO year which might be different)
func GetWeekForDate(date time.Time) int {
	year := date.Year()
	isoYear, isoWeek := date.ISOWeek()
	// If we're in December but the ISO week belongs to next year, use the last week of this year
	if date.Month() == time.December && isoYear > year {
		// Get the last day of the year
		lastDay := time.Date(year, 12, 31, 0, 0, 0, 0, date.Location())
		// Get the week of the last day
		_, lastWeek := lastDay.ISOWeek()
		return lastWeek
	}
	// If we're in January but the ISO week belongs to previous year, use week 1
	if date.Month() == time.January && isoYear < year {
		return 1
	}
	return isoWeek
}

// GetMonthName takes a time.Time and returns the name of the month in the current locale
func GetMonthName(cal kal.Calendar, t time.Time) string {
	return capitalize(cal.MonthName(t.Month()))
}

// GetMonthAbbrev takes a time.Month and returns the abbreviation in the current locale
func GetMonthAbbrev(cal kal.Calendar, month time.Month) string {
	return strings.ToLower(cal.MonthName(month)[:3])
}

// MonthNumber returns the month number given a year int and a week int.
func MonthNumber(year, week int) int {
	// Get the first Monday of the specified week
	mondayOfWeek := FirstMondayOfWeek(year, week)
	// Return the month number
	return int(mondayOfWeek.Month())
}
