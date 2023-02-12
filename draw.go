package main

import (
	"math/rand"
	"time"

	"github.com/fogleman/ln/ln"
	"github.com/signintech/gopdf"
)

// monthNumber returns the month number given a year int and a week int.
func monthNumber(year, week int) int {
	// Get the first day of the year
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	// Get the first day of the week
	firstWeekDay := firstDay.AddDate(0, 0, (week-1)*7)
	// Return the month number
	return int(firstWeekDay.Month())
}

// drawLineImage draws an image into the PDF, using only lines
func drawLineImage(pdf *gopdf.GoPdf, year, week int, x, y, width, height float64) error {

	var (
		month = monthNumber(year, week)
		paths ln.Paths
		scene = ln.Scene{}
	)

	// Thanks to github.com/fogleman for the excellent packages and examples!

	switch month % 3 {
	case 0:
		// (sphere & cube) - (cylinder & cylinder & cylinder)
		shape := ln.NewDifference(
			ln.NewIntersection(
				ln.NewSphere(ln.Vector{}, 1),
				ln.NewCube(ln.Vector{-0.8, -0.8, -0.8}, ln.Vector{0.8, 0.8, 0.8}),
			),
			ln.NewCylinder(0.4, -2, 2),
			ln.NewTransformedShape(ln.NewCylinder(0.4, -2, 2), ln.Rotate(ln.Vector{1, 0, 0}, ln.Radians(90))),
			ln.NewTransformedShape(ln.NewCylinder(0.4, -2, 2), ln.Rotate(ln.Vector{0, 1, 0}, ln.Radians(90))),
		)
		rotationAngleInDegrees := (week / 52.0) * 360.0 // rotate the csg shapes a full round through a year
		m := ln.Rotate(ln.Vector{0, 0, 1}, ln.Radians(float64(rotationAngleInDegrees)))
		scene.Add(ln.NewTransformedShape(shape, m))
		eye := ln.Vector{0, 6, 2}
		center := ln.Vector{0, 0, 0}
		up := ln.Vector{0, 0, 1}
		paths = scene.Render(eye, center, up, width, height, 20, 0.1, 100, 0.01)
	case 1:
		box := ln.Box{ln.Vector{-2, -2, -4}, ln.Vector{2, 2, 2}}
		scene.Add(ln.NewFunction(func(x, y float64) float64 {
			return -1 / (x*x + y*y)
		}, box, ln.Below))
		eye := ln.Vector{3, 0, 3}
		center := ln.Vector{1.1, 0, 0}
		up := ln.Vector{0, 0, 1}
		paths = scene.Render(eye, center, up, width, height, 50, 0.1, 100, 0.01)
	case 2:
		eye := ln.Vector{8, 8, 8}
		center := ln.Vector{0, 0, 0}
		up := ln.Vector{0, 0, 1}
		scene := ln.Scene{}
		n := 10
		for x := -n; x <= n; x++ {
			for y := -n; y <= n; y++ {
				z := rand.Float64() * 3
				v := ln.Vector{float64(x), float64(y), z}
				sphere := ln.NewOutlineSphere(eye, up, v, 0.45)
				scene.Add(sphere)
			}
		}
		fovy := 50.0
		paths = scene.Render(eye, center, up, width, height, fovy, 0.1, 100, 0.01)
	}

	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)

	px := 0.0
	py := 0.0
	scale := 1.0
	for _, path := range paths {
		for _, v := range path {
			if px != 0.0 && py != 0.0 {
				pdf.Line(x+px, y+py, x+v.X*scale, y+v.Y*scale)
			}
			px = v.X * scale
			py = v.Y * scale
		}
		px = 0.0
		py = 0.0
	}

	return nil
}
