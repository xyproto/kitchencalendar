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

var verboseLogging = env.Bool("VERBOSE")

func logVerbose(message string) {
	if verboseLogging {
		fmt.Println(message)
	}
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

	if req.WeeksSpan <= 0 {
		req.WeeksSpan = 1 // Default to 1 week if WeeksSpan is not specified or invalid
	}

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	start := fromDate
	for firstIteration := true; firstIteration || start.Before(toDate); firstIteration = false {
		end := start.AddDate(0, 0, 7*req.WeeksSpan-1)
		if end.After(toDate) {
			end = toDate
		}

		year, week := start.ISOWeek()
		logVerbose(fmt.Sprintf("Generating PDF for week %d of year %d", week, year))

		pdfBytes, err := kc.GeneratePDF(year, week, req.Names, req.Drawing)
		if err != nil {
			http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
			logVerbose(fmt.Sprintf("Error generating PDF: %v", err))
			return
		}

		fileName := fmt.Sprintf("calendar_%d-%d.pdf", year, week)
		fw, err := zipWriter.Create(fileName)
		if err != nil {
			http.Error(w, "Failed to create zip entry", http.StatusInternalServerError)
			logVerbose(fmt.Sprintf("Error creating zip entry: %v", err))
			return
		}
		_, err = fw.Write(pdfBytes)
		if err != nil {
			http.Error(w, "Failed to write to zip", http.StatusInternalServerError)
			logVerbose(fmt.Sprintf("Error writing to zip: %v", err))
			return
		}

		// Skip a week, since each PDF has two weeks when iterating
		start = end.AddDate(0, 0, 7)
	}

	if err := zipWriter.Close(); err != nil {
		http.Error(w, "Failed to close zip writer", http.StatusInternalServerError)
		logVerbose(fmt.Sprintf("Error closing zip writer: %v", err))
		return
	}

	logVerbose("PDF generation and zipping completed successfully")
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=calendars.zip")
	if _, err := io.Copy(w, buffer); err != nil {
		logVerbose(fmt.Sprintf("Error sending zip file: %v", err))
		return
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
