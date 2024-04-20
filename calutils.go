package kitchencalendar

import (
	"strings"
	"time"

	"github.com/xyproto/kal"
)

// FirstMondayOfWeek finds the first monday of the week, given a year and a week number
func FirstMondayOfWeek(year, week int) time.Time {
	// Create a time object for the given year
	t := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// Calculate the number of days to add to the time object
	// to get the first Monday of the given week number
	daysToAdd := (week - 1) * 7
	t = t.AddDate(0, 0, daysToAdd)
	// If the day of the week is not Monday,
	// add the number of days to get to the next Monday
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

// FirstSundayOfWeek finds the first monday of the week, given a year and a week number
func FirstSundayOfWeek(year, week int) time.Time {
	// Create a time object for the given year
	t := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// Calculate the number of days to add to the time object
	// to get the first Monday of the given week number
	daysToAdd := (week - 1) * 7
	t = t.AddDate(0, 0, daysToAdd)
	// If the day of the week is not Sunday,
	// add the number of days to get to the next Sunday
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
	}
	return t
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
	_, week := now.ISOWeek()
	// Return the week number
	return week
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
	// Get the first day of the year
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	// Get the first day of the week
	firstWeekDay := firstDay.AddDate(0, 0, (week-1)*7)
	// Return the month number
	return int(firstWeekDay.Month())
}
