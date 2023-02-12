package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/signintech/gopdf"
	"github.com/xyproto/env/v2"
	"github.com/xyproto/kal"
)

const versionString = "KitchenCalendar 0.0.1"

var paperSize = env.Str("PAPERSIZE", "A4")

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

// firstSundayOfWeek finds the first monday of the week, given a year and a week number
func firstSundayOfWeek(year, week int) time.Time {
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

// firstSundayAfter finds the first Sunday after the given date
func firstSundayAfter(date time.Time) time.Time {
	// Get the day of the week for the given date
	dayOfWeek := date.Weekday()
	// Calculate the number of days to add to the given date to get the first Sunday
	daysToAdd := 7 - int(dayOfWeek)
	// Add the calculated number of days to the given date and return it
	return date.AddDate(0, 0, daysToAdd)
}

// FirstSaturdayAfter takes a time.Time and returns the first Saturday after the given date as a time.Time
func firstSaturdayAfter(date time.Time) time.Time {
	// Get the day of the week for the given date
	dayOfWeek := date.Weekday()
	// Calculate the number of days until the next Saturday
	daysUntilSaturday := 6 - int(dayOfWeek)
	// Add the number of days until the next Saturday to the given date
	return date.AddDate(0, 0, daysUntilSaturday)
}

// generateTitle generates the main title of the calendar
func generateTitle(cal kal.Calendar, year, week int) string {
	mondayTime := firstMondayOfWeek(year, week)
	monthName1 := getMonthName(cal, mondayTime)
	week++
	mondayTime = firstMondayOfWeek(year, week)
	monthName2 := getMonthName(cal, mondayTime)
	if monthName1 == monthName2 {
		return fmt.Sprintf("%s %d", monthName1, year)
	}
	return fmt.Sprintf("%s - %s %d", monthName1, monthName2, year)
}

// generateWeekHeaderLeft creates the header for the right side of the week table
// on the format: from date -> to date
func generateWeekHeaderRight(cal kal.Calendar, year, week int) string {
	mondayTime := firstMondayOfWeek(year, week)
	sundayTime := firstSundayAfter(mondayTime)
	return fmt.Sprintf("%s -> %s", formatDate(cal, mondayTime), formatDate(cal, sundayTime))
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

// draw a week into the PDF
func drawWeek(pdf *gopdf.GoPdf, cal kal.Calendar, year, week int, x, y *float64, xRight, width float64, names []string) error {
	tableHeight := 300.0

	// Draw the left vertical lines of the table
	pdf.Line(*x, *y+20, *x, *y+tableHeight+36)

	// Draw the right vertical lines of the table
	pdf.Line(*x+width, *y+20, *x+width, *y+tableHeight+36)

	// Generate the titles for this week
	headerLeft := generateWeekHeaderLeft(year, week)
	headerRight := generateWeekHeaderRight(cal, year, week)

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
	sundayTime := firstSundayAfter(mondayTime)

	// Draw the week names and vertical lines for the 1st week
	originalX := *x
	*x += 70
	err := iterateDays(mondayTime, sundayTime, func(t time.Time) error {
		text := dayAndDate(cal, t)

		fontName := "regular"
		if isRedDay := kal.RedDay(cal, t); t.Weekday() == time.Sunday || isRedDay { // Red day
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

func drawImage(pdf *gopdf.GoPdf, year, week int, x, y, w, h float64) {
	r := rand.New(rand.NewSource(int64(year)*256 + int64(week)))
	maxLineWidth := 2.0
	pdf.SetLineWidth(r.Float64() * maxLineWidth)
	px1 := float64(r.Intn(int(w)))
	py1 := float64(r.Intn(int(h)))
	for i := 0; i < 20; i++ {
		px2 := float64(r.Intn(int(w)))
		py2 := float64(r.Intn(int(h)))
		pdf.Line(x+px1, y+py1, x+px2, y+py2)
		px1 = px2
		py1 = py2
	}
	maxLineWidth = 3.0
	for i := 0; i < 5; i++ {
		px1 := float64(r.Intn(int(w)))
		py1 := float64(r.Intn(int(h)))
		px2 := float64(r.Intn(int(w)))
		py2 := float64(r.Intn(int(h)))
		pdf.SetLineWidth(r.Float64() * maxLineWidth)
		pdf.Oval(x+px1, y+py1, x+px2, y+py2)
	}
}

// capitalize makes changes the first rune of a string to be in uppercase
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	firstRune := unicode.ToUpper(runes[0])
	return string(append([]rune{firstRune}, runes[1:]...))
}

// getMonthName takes a time.Time and returns the name of the month in the current locale
func getMonthName(cal kal.Calendar, t time.Time) string {
	return capitalize(cal.MonthName(t.Month()))
}

// getMonthAbbrev takes a time.Month and returns the abbreviation in the current locale
func getMonthAbbrev(cal kal.Calendar, month time.Month) string {
	return strings.ToLower(cal.MonthName(month)[:3])
}

func main() {
	outputFilename := flag.String("o", "calendar.pdf", "an output PDF filename")
	yearFlag := flag.Int("year", getCurrentYear(), "the year")
	weekFlag := flag.Int("week", getCurrentWeek(), "the week number")
	nameString := flag.String("names", "Bob,Alice,Mallory,Judy", "names used in the calendar")
	drawingFlag := flag.Bool("drawing", true, "include a drawing for each year and week in the top right corner")
	verbose := flag.Bool("V", true, "verbose output")

	flag.Parse()

	year := *yearFlag
	week := *weekFlag
	names := strings.Split(*nameString, ",")

	cal, err := newCalendar()
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

	y := 60.0
	x := 35.0
	xRight := 460.0
	width := 538.0

	// Draw the month and year title
	title := generateTitle(cal, year, week)
	if err := write(&pdf, x, y, title, "bold", 24); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Draw a little logo for this year and week in the top right, by using random lines
	if *drawingFlag {
		drawImage(&pdf, year, week, xRight-50, y-40, 170, 70)
	}

	// Set the line width for the weeks and tables that will now be drawn
	pdf.SetLineWidth(0.5)

	// Draw the first week
	y += 50
	if err := drawWeek(&pdf, cal, year, week, &x, &y, xRight, width, names); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	week++

	// Draw the second week
	y += 20
	if err := drawWeek(&pdf, cal, year, week, &x, &y, xRight, width, names); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if *verbose {
		fmt.Printf("Writing to %s... ", *outputFilename)
	}

	pdf.WritePdf(*outputFilename)

	if *verbose {
		fmt.Println("done")
	}

}