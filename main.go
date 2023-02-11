package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xyproto/env"
)

const versionString = "KitchenCalendar 0.0.1"

const verbose = true

var paperSize = env.Str("PAPERSIZE", "A4")

func init() {
	fmt.Println(versionString)
}

func main() {
	outputFilenameFlag := flag.String("o", "calendar.pdf", "an output PDF filename")
	verboseFlag := flag.Bool("V", true, "verbose output")
	flag.Parse()

	filename := *outputFilenameFlag
	verbose := *verboseFlag

	timestamp := time.Now().Format("2006-01-02")

	pdf := gofpdf.New("P", "mm", paperSize, "")
	pdf.SetTopMargin(30)
	topLeftText := "1/1"
	topRightText := timestamp + ", " + "BLABLABLA"
	pdf.SetHeaderFunc(func() {
		pdf.SetY(5)
		pdf.SetFont("Helvetica", "", 6)
		pdf.CellFormat(80, 0, topLeftText, "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 0, topRightText, "", 0, "R", false, 0, "")
	})
	pdf.AddPage()
	pdf.SetY(20)
	lines := strings.Split("this\nis\ntext", "\n")
	pdf.SetFont("Courier", "B", 12)
	pdf.Write(5, lines[0]+"\n")
	pdf.SetFont("Courier", "", 12)
	pdf.Write(5, strings.Join(lines[1:len(lines)-1], "\n"))
	pdf.SetFont("Courier", "B", 12)
	pdf.Write(5, "\n"+lines[len(lines)-1])

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s already exists\n", filename)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Writing %s... ", filename)
	}
	if err := pdf.OutputFileAndClose(filename); err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(1)
	}
	if verbose {
		fmt.Println("done")
	}
}
