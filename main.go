package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/signintech/gopdf"
	"github.com/xyproto/env"
)

const versionString = "KitchenCalendar 0.0.1"

var paperSize = env.Str("PAPERSIZE", "A4")

func init() {
	fmt.Println(versionString)
}

func main() {
	outputFilename := flag.String("o", "calendar.pdf", "an output PDF filename")
	verbose := flag.Bool("V", true, "verbose output")
	flag.Parse()

	//timestamp := time.Now().Format("2006-01-02")

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
	err := pdf.AddTTFFont("wts11", "ttf/nunito/Nunito-Regular.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}

	err = pdf.SetFont("wts11", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}
	pdf.Cell(nil, "asdf")
	pdf.WritePdf(*outputFilename)

	if *verbose {
		fmt.Println("done")
	}

}
