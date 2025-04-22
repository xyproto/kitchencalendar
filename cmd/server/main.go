//go:build nb_NO || en_US

package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/xyproto/env/v2"
	kc "github.com/xyproto/kitchencalendar"
)

type CalendarRequest struct {
	FromDate  string   `json:"fromDate"`
	ToDate    string   `json:"toDate"`
	Names     []string `json:"names"`
	Drawing   bool     `json:"drawing"`
	WeeksSpan int      `json:"weeksSpan"` // 1 or 2 weeks per PDF
}

const (
	daysPerWeek      = 7
	defaultWeeksSpan = 2
)

var verboseLogging = env.Bool("VERBOSE")

func logVerbose(message string) {
	if verboseLogging {
		fmt.Println(message)
	}
}

func generateCalendars(req CalendarRequest, fromDate, toDate time.Time) ([]byte, error) {
	if req.WeeksSpan <= 0 {
		req.WeeksSpan = defaultWeeksSpan // Default to defaultWeeksSpan if WeeksSpan is not specified or invalid
	}

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	start := fromDate
	for firstIteration := true; firstIteration || start.Before(toDate); firstIteration = false {
		expectedStart := start
		end := start.AddDate(0, 0, daysPerWeek*req.WeeksSpan-1)
		if end.After(toDate) {
			end = toDate
		}

		year := start.Year()
		week := kc.GetWeekForDate(start)
		logVerbose(fmt.Sprintf("Generating PDF for weeks %d-%d of year %d", week, week+req.WeeksSpan-1, year))

		pdfBytes, err := kc.GeneratePDF(year, week, req.Names, req.Drawing)
		if err != nil {
			return nil, fmt.Errorf("failed to generate PDF: %v", err)
		}

		fileName := fmt.Sprintf("calendar_%d-%d.pdf", year, week)
		fw, err := zipWriter.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create zip entry: %v", err)
		}
		if _, err := fw.Write(pdfBytes); err != nil {
			return nil, fmt.Errorf("failed to write to zip: %v", err)
		}

		// Move to the next set of weeks, avoiding overlap
		start = start.AddDate(0, 0, daysPerWeek*req.WeeksSpan)

		// Check if the iteration is progressing as expected
		if start.Before(expectedStart.AddDate(0, 0, daysPerWeek*req.WeeksSpan)) {
			fmt.Fprintf(os.Stderr, "Warning: Iteration did not progress as expected. Current start: %v, Expected start: %v\n", start, expectedStart.AddDate(0, 0, daysPerWeek*req.WeeksSpan))
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %v", err)
	}

	logVerbose("PDF generation and zipping completed successfully")
	return buffer.Bytes(), nil
}

func handleCreateCalendar(w http.ResponseWriter, r *http.Request) {
	logVerbose("Received request to create calendar")

	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		logVerbose("Error: Method not supported")
		return
	}

	var req CalendarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		logVerbose(fmt.Sprintf("Error decoding request: %v", err))
		return
	}
	logVerbose(fmt.Sprintf("Request decoded successfully: %+v", req))

	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		http.Error(w, "Invalid from date", http.StatusBadRequest)
		logVerbose(fmt.Sprintf("Error parsing from date: %v", err))
		return
	}

	toDate, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		http.Error(w, "Invalid to date", http.StatusBadRequest)
		logVerbose(fmt.Sprintf("Error parsing to date: %v", err))
		return
	}

	zipData, err := generateCalendars(req, fromDate, toDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logVerbose(fmt.Sprintf("Error generating calendars: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=calendars.zip")
	if _, err := io.Copy(w, bytes.NewReader(zipData)); err != nil {
		logVerbose(fmt.Sprintf("Error sending zip file: %v", err))
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/style.css")
	})

	http.HandleFunc("/createcalendar", handleCreateCalendar)

	fmt.Println("Serving on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
