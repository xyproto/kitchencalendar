package main

import (
	"math"
	"math/rand"

	"github.com/fogleman/ln/ln"
	"github.com/signintech/gopdf"
)

func Normalize(values []float64, a, b float64) []float64 {
	result := make([]float64, len(values))
	lo := values[0]
	hi := values[0]
	for _, x := range values {
		lo = math.Min(lo, x)
		hi = math.Max(hi, x)
	}
	for i, x := range values {
		p := (x - lo) / (hi - lo)
		result[i] = a + p*(b-a)
	}
	return result
}

func LowPass(values []float64, alpha float64) []float64 {
	result := make([]float64, len(values))
	var y float64
	for i, x := range values {
		y -= alpha * (y - x)
		result[i] = y
	}
	return result
}

func LowPassNoise(n int, alpha float64, iterations int) []float64 {
	result := make([]float64, n)
	for i := range result {
		result[i] = rand.Float64()
	}
	for i := 0; i < iterations; i++ {
		result = LowPass(result, alpha)
	}
	result = Normalize(result, -1, 1)
	return result
}

type Tree struct {
	ln.Shape
	V0, V1 ln.Vector
}

func (t *Tree) Paths() ln.Paths {
	paths := t.Shape.Paths()
	for i := 0; i < 128; i++ {
		p := math.Pow(rand.Float64(), 1.5)*0.5 + 0.5
		c := t.V0.Add(t.V1.Sub(t.V0).MulScalar(p))
		a := rand.Float64() * 2 * math.Pi
		l := (1 - p) * 8
		d := ln.Vector{math.Cos(a), math.Sin(a), -2.75}.Normalize()
		e := c.Add(d.MulScalar(l))
		paths = append(paths, ln.Path{c, e})
	}
	return paths
}

func drawLineImage(pdf *gopdf.GoPdf, year, week int, x, y, width, height float64) error {

	var (
		sheetNumber = int(math.Round(float64(week) * 0.5)) // which sheet number
		paths       ln.Paths
		scene       = ln.Scene{}
	)

	// Thanks to github.com/fogleman for the excellent packages and examples!

	if sheetNumber%2 == 0 {
		rotationAngleInDegrees := (week / 52.0) * 360.0 // rotate the csg shapes a full round through a year
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
		m := ln.Rotate(ln.Vector{0, 0, 1}, ln.Radians(float64(rotationAngleInDegrees)))
		scene.Add(ln.NewTransformedShape(shape, m))
		eye := ln.Vector{0, 6, 2}
		center := ln.Vector{0, 0, 0}
		up := ln.Vector{0, 0, 1}
		paths = scene.Render(eye, center, up, width, height, 20, 0.1, 100, 0.01)
	} else {

		box := ln.Box{ln.Vector{-2, -2, -4}, ln.Vector{2, 2, 2}}
		scene.Add(ln.NewFunction(func(x, y float64) float64 {
			return -1 / (x*x + y*y)
			// return math.Cos(x*y) * (x*x - y*y)
		}, box, ln.Below))
		eye := ln.Vector{3, 0, 3}
		center := ln.Vector{1.1, 0, 0}
		up := ln.Vector{0, 0, 1}
		paths = scene.Render(eye, center, up, width, height, 50, 0.1, 100, 0.01)

	}

	//pdf.SetStrokeColor(80, 80, 80)
	//pdf.SetStrokeColor(180, 180, 180)
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
