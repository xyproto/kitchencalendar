package kitchencalendar

import (
	_ "embed"
	"errors"
	"path/filepath"

	"fmt"
	"os"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/xyproto/env/v2"
	"github.com/xyproto/kal"
)

var paperSize = env.Str("PAPERSIZE", "A4")

//go:embed ttf/nunito/Nunito-Regular.ttf
var nunitoRegularData []byte

//go:embed ttf/nunito/Nunito-Bold.ttf
var nunitoBoldData []byte

// generateTitle generates the main title of the calendar
func generateTitle(cal kal.Calendar, year, week int) string {
	mondayTime := FirstMondayOfWeek(year, week)
	monthName1 := GetMonthName(cal, mondayTime)
	week++
	mondayTime = FirstMondayOfWeek(year, week)
	monthName2 := GetMonthName(cal, mondayTime)
	if monthName1 == monthName2 {
		return fmt.Sprintf("%s %d", monthName1, year)
	}
	return fmt.Sprintf("%s - %s %d", monthName1, monthName2, year)
}

// generateWeekHeaderLeft creates the header for the right side of the week table
// on the format: from date -> to date
func generateWeekHeaderRight(cal kal.Calendar, year, week int) string {
	mondayTime := FirstMondayOfWeek(year, week)
	sundayTime := FirstSundayAfter(mondayTime)
	return fmt.Sprintf("%s -> %s", FormatDate(cal, mondayTime), FormatDate(cal, sundayTime))
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
	headerLeft := WeekString(week)
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
	mondayTime := FirstMondayOfWeek(year, week)
	sundayTime := FirstSundayAfter(mondayTime)

	// Draw the week names and vertical lines for the 1st week
	originalX := *x
	*x += 70
	cellWidth := width / 8.0
	i := 1
	err := IterateDays(mondayTime, sundayTime, func(t time.Time) error {
		text := DayAndDate(cal, t)

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

	if len(names) == 0 {
		return errors.New("the given slice of names is empty")
	}

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

func GeneratePDF(year, week int, names []string, drawing bool) ([]byte, error) {
	cal, err := NewCalendar()
	if err != nil {
		return []byte{}, err
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
		if err := os.WriteFile(nunitoRegularFilename, nunitoRegularData, 0o664); err != nil {
			return []byte{}, fmt.Errorf("could not write to %s: %w", nunitoRegularFilename, err)
		}
	}
	if !exists(nunitoRegularFilename) {
		return []byte{}, fmt.Errorf("could not write to %s", nunitoRegularFilename)
	}
	defer os.Remove(nunitoRegularFilename)

	nunitoBoldFilename := filepath.Join(tempdir, "Nunito-Bold.ttf")
	if !exists(nunitoBoldFilename) {
		if err := os.WriteFile(nunitoBoldFilename, nunitoBoldData, 0o664); err != nil {
			return []byte{}, fmt.Errorf("could not write to %s: %w", nunitoBoldFilename, err)
		}
	}
	if !exists(nunitoBoldFilename) {
		return []byte{}, fmt.Errorf("could not write to %s", nunitoBoldFilename)
	}
	defer os.Remove(nunitoBoldFilename)

	if err := pdf.AddTTFFont("regular", nunitoRegularFilename); err != nil {
		return []byte{}, err
	}

	if err := pdf.AddTTFFont("bold", nunitoBoldFilename); err != nil {
		return []byte{}, err
	}

	y := 35.0
	x := 35.0
	width := 538.0

	// Draw the month and year title
	title := generateTitle(cal, year, week)
	if err := write(&pdf, x, y, title, "bold", 24); err != nil {
		return []byte{}, err
	}

	if drawing {
		DrawLineImage(&pdf, year, week, width-40, y-10, 70, 70)
	}

	// Set the line width for the weeks and tables
	pdf.SetLineWidth(1.0)

	// Draw the first week
	y += 75
	if err := drawWeek(&pdf, cal, year, week, &x, &y, width, names); err != nil {
		return []byte{}, err
	}

	week++

	// Draw the second week
	y += 20
	if err := drawWeek(&pdf, cal, year, week, &x, &y, width, names); err != nil {
		return []byte{}, err
	}

	return pdf.GetBytesPdf(), nil
}
