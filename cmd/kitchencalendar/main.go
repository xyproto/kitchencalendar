package main

import (
	_ "embed"

	"flag"
	"fmt"
	"os"
	"strings"

	kc "github.com/xyproto/kitchencalendar"
)

func main() {
	outputFilename := flag.String("o", "", "an output PDF filename")
	yearFlag := flag.Int("year", kc.GetCurrentYear(), "the year")
	weekFlag := flag.Int("week", kc.GetCurrentWeek(), "the week number")
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

	pdfBytes, err := kc.GeneratePDF(year, week, names, *drawingFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if *verbose {
		fmt.Printf("Writing %s... ", filename)
	}

	if err := os.WriteFile(filename, pdfBytes, 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if *verbose {
		fmt.Println("done")
	}

}
