// This file provides a point abstraction that is analogous to image.Point but
// based on floating-point rather than integer coordinates.

package xmorph

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// A Point is a float64-valued (x, y) coordinate.  The axes increase right and
// down.
type Point struct {
	X float64
	Y float64
}

// Add adds one Point to another to produce a third Point.
func (p Point) Add(q Point) Point {
	return Point{
		X: p.X + q.X,
		Y: p.Y + q.Y,
	}
}

// Sub subtracts one Point from another to produce a third Point.
func (p Point) Sub(q Point) Point {
	return Point{
		X: p.X - q.X,
		Y: p.Y - q.Y,
	}
}

// Mul multiplies a Point by a scalar to produce a new Point.
func (p Point) Mul(k float64) Point {
	return Point{
		X: p.X * k,
		Y: p.Y * k,
	}
}

// Div divides a Point by a scalar to produce a new Point.
func (p Point) Div(k float64) Point {
	return Point{
		X: p.X / k,
		Y: p.Y / k,
	}
}

// Eq reports whether two Points are equal within a given tolerance.  The
// tolerance is applied separately to the x and y coordinates; both must be
// within tolerance for the function to return true.
func (p Point) Eq(q Point, tol float64) bool {
	if tol < 0.0 {
		panic("Eq tolerance must be non-negative")
	}
	dx := math.Abs(p.X - q.X)
	dy := math.Abs(p.Y - q.Y)
	return dx <= tol && dy <= tol
}

// formatString performs most of the work for Format.  The difference is that
// it returns a string, which Format sends to the correct receiver.
func (p Point) formatString(st fmt.State, verb rune) string {
	type RawPoint Point // Fresh method set
	switch verb {
	case 'v':
		if st.Flag('#') {
			str := fmt.Sprintf("%#v", RawPoint(p))
			str = strings.Replace(str, "RawPoint", "Point", 1)
			return str
		} else {
			return fmt.Sprintf("[%v, %v]", p.X, p.Y)
		}
	case 'b', 'e', 'E', 'f', 'F', 'g', 'G', 'x', 'X':
		// Propagate all known flags when provided.
		fstr := "%"
		fl := make([]rune, 0, 5)
		for _, c := range "+-# 0" {
			if st.Flag(int(c)) {
				fl = append(fl, c)
			}
		}
		fstr += string(fl)

		// Propagate the width and precision when provided.
		if wd, ok := st.Width(); ok {
			fstr += strconv.Itoa(wd)
		}
		if prec, ok := st.Precision(); ok {
			fstr += "." + strconv.Itoa(prec)
		}

		// Retain the verb as is.
		fstr += string(verb)

		// Apply the format to X and Y.
		return fmt.Sprintf("["+fstr+", "+fstr+"]", p.X, p.Y)
	}
	return "%![invalid Point format]"
}

// Format applies standard numeric formatting to a Point's coordinates when
// outputting it via fmt.Printf et al.
func (p Point) Format(st fmt.State, verb rune) {
	fmt.Fprintf(st, p.formatString(st, verb))
}
