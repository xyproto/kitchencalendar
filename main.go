package main

import (
	_ "embed"
	"path/filepath"

	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/xyproto/env/v2"
	"github.com/xyproto/kal"
)

const versionString = "KitchenCalendar 0.2.0"

//go:embed ttf/nunito/Nunito-Regular.ttf
var nunitoRegularData []byte

//go:embed ttf/nunito/Nunito-Bold.ttf
var nunitoBoldData []byte

var paperSize = env.Str("PAPERSIZE", "A4")

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

// draw a week into the PDF
func drawWeek(pdf *gopdf.GoPdf, cal kal.Calendar, year, week int, x, y *float64, width float64, names []string) error {
	tableHeight := 300.0

	// Draw the left vertical lines of the table
	pdf.Line(*x, *y+20, *x, *y+tableHeight+37.2)

	// Draw the right vertical lines of the table
	pdf.Line(*x+width, *y+20, *x+width, *y+tableHeight+37.2)

	// Generate the titles for this week
	headerLeft := weekString(week)
	headerRight := generateWeekHeaderRight(cal, year, week)

	// Draw the header for the 1st week
	if err := write(pdf, *x, *y, headerLeft, "bold", 14); err != nil {
		return err
	}
	approxHeaderRightWidth := float64(len(headerRight)) * 4.9
	if err := write(pdf, width-approxHeaderRightWidth, *y, headerRight, "regular", 14); err != nil {
		return err
	}

	// Draw a horizontal line
	*y += 20
	pdf.Line(*x-0.2, *y, *x+width+0.2, *y)

	// Find monday and sunday
	mondayTime := firstMondayOfWeek(year, week)
	sundayTime := firstSundayAfter(mondayTime)

	// Draw the week names and vertical lines for the 1st week
	originalX := *x
	*x += 70
	cellWidth := width / 8.0
	i := 1
	err := iterateDays(mondayTime, sundayTime, func(t time.Time) error {
		text := dayAndDate(cal, t)

		fontName := "regular"
		if isRedDay := kal.RedDay(cal, t); t.Weekday() == time.Sunday || isRedDay { // Red day
			fontName = "bold"
		}

		fontSize := 11
		if err := write(pdf, originalX+float64(i)*cellWidth+2, *y, text, fontName, fontSize); err != nil {
			return err
		}
		// Draw the vertical line
		pdf.Line(originalX+float64(i)*cellWidth, *y, originalX+float64(i)*cellWidth, *y+tableHeight+17.3)
		i++
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

func main() {
	outputFilename := flag.String("o", "", "an output PDF filename")
	yearFlag := flag.Int("year", getCurrentYear(), "the year")
	weekFlag := flag.Int("week", getCurrentWeek(), "the week number")
	nameString := flag.String("names", "Bob,Alice,Mallory,Judy", "names used in the calendar")
	drawingFlag := flag.Bool("drawing", true, "include a drawing for each year and week in the top right corner")
	verbose := flag.Bool("V", true, "verbose output")

	flag.Parse()

	year := *yearFlag
	week := *weekFlag
	names := strings.Split(*nameString, ",")

	filename := ""
	if *outputFilename == "" {
		filename = fmt.Sprintf("calendar_w%d_%d.pdf", week, year)
	} else {
		filename = *outputFilename
	}

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

	tempdir := env.Str("TMPDIR", "/tmp")
	nunitoRegularFilename := filepath.Join(tempdir, "Nunito-Regular.ttf")
	if !exists(nunitoRegularFilename) {
		os.WriteFile(nunitoRegularFilename, nunitoRegularData, 0o664)
	}
	if !exists(nunitoRegularFilename) {
		err := fmt.Errorf("could not write to %s", nunitoRegularFilename)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	nunitoBoldFilename := filepath.Join(tempdir, "Nunito-Bold.ttf")
	if !exists(nunitoBoldFilename) {
		os.WriteFile(nunitoBoldFilename, nunitoBoldData, 0o664)
	}

	if !exists(nunitoBoldFilename) {
		err := fmt.Errorf("could not write to %s", nunitoBoldFilename)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := pdf.AddTTFFont("regular", nunitoRegularFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := pdf.AddTTFFont("bold", nunitoBoldFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	y := 35.0
	x := 35.0
	width := 538.0

	// Draw the month and year title
	title := generateTitle(cal, year, week)
	if err := write(&pdf, x, y, title, "bold", 24); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if *drawingFlag {
		drawLineImage(&pdf, year, week, width-40, y-10, 70, 70)
	}

	// Set the line width for the weeks and tables that will now be drawn
	pdf.SetLineWidth(0.1)

	// Draw the first week
	y += 75
	if err := drawWeek(&pdf, cal, year, week, &x, &y, width, names); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	week++

	// Draw the second week
	y += 20
	if err := drawWeek(&pdf, cal, year, week, &x, &y, width, names); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if *verbose {
		fmt.Printf("Writing to %s... ", filename)
	}
	pdf.WritePdf(filename)
	if *verbose {
		fmt.Println("done")
	}

}
