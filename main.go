package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/xyproto/env"
	"github.com/xyproto/kal"
)

const (
	versionString = "KitchenCalendar 0.0.1"
)

var paperSize = env.Str("PAPERSIZE", "A4")

func init() {
	fmt.Println(versionString)
}

// firstMondayOfWeek finds the first monday of the week, given a year and a week number
func firstMondayOfWeek(year, week int) time.Time {
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

// findFirstSunday finds the first sunday after the given date
func findFirstSunday(date time.Time) time.Time {
	// Get the day of the week for the given date
	dayOfWeek := date.Weekday()
	// Calculate the number of days to add to the given date to get the first Sunday
	daysToAdd := 7 - int(dayOfWeek)
	// Add the calculated number of days to the given date
	firstSunday := date.AddDate(0, 0, daysToAdd)
	// Return the first Sunday as a string
	return firstSunday
}

// getMonthNameInNorwegian takes a time.Time and returns the name of the month in Norwegian
func getMonthNameInNorwegian(t time.Time) string {
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

// generateTitle generates the main title of the calendar
func generateTitle(year, week int) string {
	mondayTime := firstMondayOfWeek(year, week)
	monthName1 := getMonthNameInNorwegian(mondayTime)
	week++
	mondayTime = firstMondayOfWeek(year, week)
	monthName2 := getMonthNameInNorwegian(mondayTime)
	if monthName1 == monthName2 {
		return fmt.Sprintf("%s %d", monthName1, year)
	}
	return fmt.Sprintf("%s - %s %d", monthName1, monthName2, year)
}

// generateWeekHeaderLeft creates the header for the left side of the week table
// on the format "Uke N"
func generateWeekHeaderLeft(year, week int) string {
	return fmt.Sprintf("Uke %d", week)
}

// generateWeekHeaderLeft creates the header for the right side of the week table
// on the format: from date -> to date
func generateWeekHeaderRight(year, week int) string {
	mondayTime := firstMondayOfWeek(year, week)
	sundayTime := findFirstSunday(mondayTime)
	return fmt.Sprintf("%s -> %s", formatDate(mondayTime), formatDate(sundayTime))
}

func write(pdf *gopdf.GoPdf, x, y float64, text string, fontName string, fontSize int) error {
	if err := pdf.SetFont(fontName, "", fontSize); err != nil {
		return err
	}
	pdf.SetXY(x, y)
	pdf.Cell(nil, text)
	return nil
}

// iterateDays iterates over days from startDay to endDay (inclusive) and calls f for each day
func iterateDays(startDay, endDay time.Time, f func(time.Time) error) error {
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

// draw a week into the PDF
func drawWeek(pdf *gopdf.GoPdf, calendar *kal.Calendar, year, week int, x, y *float64, xRight, width float64, names []string) error {
	tableHeight := 300.0

	// Draw the left vertical lines of the table
	pdf.Line(*x, *y+20, *x, *y+tableHeight+36)

	// Draw the right vertical lines of the table
	pdf.Line(*x+width, *y+20, *x+width, *y+tableHeight+36)

	// Generate the titles for this week
	headerLeft := generateWeekHeaderLeft(year, week)
	headerRight := generateWeekHeaderRight(year, week)

	// Draw the header for the 1st week
	if err := write(pdf, *x, *y, headerLeft, "bold", 14); err != nil {
		return err
	}
	if err := write(pdf, xRight, *y, headerRight, "regular", 14); err != nil {
		return err
	}

	// Draw a horizontal line
	*y += 20
	pdf.Line(*x, *y, *x+width, *y)

	// Find monday and sunday
	mondayTime := firstMondayOfWeek(year, week)
	sundayTime := findFirstSunday(mondayTime)

	// Draw the week names and vertical lines for the 1st week
	originalX := *x
	*x += 70
	err := iterateDays(mondayTime, sundayTime, func(t time.Time) error {
		text := dayAndDate(t)

		fontName := "regular"
		if isRedDay := kal.RedDay(*calendar, t); t.Weekday() == time.Sunday || isRedDay { // Red day
			fontName = "bold"
		}

		fontSize := 11
		if err := write(pdf, *x, *y, text, fontName, fontSize); err != nil {
			return err
		}
		pdf.Line(*x-2, *y, *x-2, *y+tableHeight+17)
		*x += float64(len(text)) * 6.5
		return nil
	})
	if err != nil {
		return err
	}
	*x = originalX

	// Draw a horizontal line
	*y += 15
	pdf.Line(*x, *y, *x+width, *y)

	nameHeight := tableHeight / float64(len(names))

	// Draw the names of the people that should use this calendar, with horizontal lines
	*y += 2
	for _, text := range names {
		// Draw the names
		fontName := "regular"
		fontSize := 12
		if err := write(pdf, *x+3, *y+1, text, fontName, fontSize); err != nil {
			return err
		}
		*y += nameHeight
		pdf.Line(*x, *y, *x+width, *y)
	}

	return nil
}

// GetCurrentYear returns the current year as an int
func getCurrentYear() int {
	return time.Now().Year()
}

// getCurrentWeek returns the current week number as an int
func getCurrentWeek() int {
	// Get the current time
	now := time.Now()
	// Get the ISO year and week number
	_, week := now.ISOWeek()
	// Return the week number
	return week
}

func main() {
	outputFilename := flag.String("o", "calendar.pdf", "an output PDF filename")
	yearFlag := flag.Int("year", getCurrentYear(), "the year")
	weekFlag := flag.Int("week", getCurrentWeek(), "the week number")
	nameString := flag.String("names", "Vilde,Synne,Aria,Alexander", "names used in the calendar")
	verbose := flag.Bool("V", true, "verbose output")

	flag.Parse()

	year := *yearFlag
	week := *weekFlag
	names := strings.Split(*nameString, ",")

	calendar, err := kal.NewCalendar("nb_NO", true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Got all needed information, generate and output the PDF

	pdf := gopdf.GoPdf{}

	// Initialize and use a config struct
	var c gopdf.Config
	switch strings.TrimSpace(paperSize) {
	case "letter":
		c.PageSize = *gopdf.PageSizeLetter
	default:
		c.PageSize = *gopdf.PageSizeA4
	}
	pdf.Start(c)

	pdf.AddPage()

	if err := pdf.AddTTFFont("regular", "ttf/nunito/Nunito-Regular.ttf"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := pdf.AddTTFFont("bold", "ttf/nunito/Nunito-Bold.ttf"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	pdf.SetLineWidth(0.5)

	y := 60.0
	x := 35.0
	xRight := 460.0
	width := 538.0

	// Draw the month and year title
	title := generateTitle(year, week)
	if err := write(&pdf, x, y, title, "bold", 24); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Draw the first week
	y += 50
	if err := drawWeek(&pdf, &calendar, year, week, &x, &y, xRight, width, names); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	week++

	// Draw the second week
	y += 20
	if err := drawWeek(&pdf, &calendar, year, week, &x, &y, xRight, width, names); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	pdf.WritePdf(*outputFilename)

	if *verbose {
		fmt.Println("done")
	}

}
