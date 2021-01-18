// The functions defined in this file ensure the xmorph package's Point
// operations work as expected.

package xmorph

import (
	"fmt"
	"math/rand"
	"testing"
)

// TestAdd tests that two points can be added.
func TestAdd(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	for i := 0; i < 100; i++ {
		pt1 := Point{
			X: rng.Float64()*21 - 11,
			Y: rng.Float64()*21 - 11,
		}
		pt2 := Point{
			X: rng.Float64()*21 - 11,
			Y: rng.Float64()*21 - 11,
		}
		exp := Point{
			X: pt1.X + pt2.X,
			Y: pt1.Y + pt2.Y,
		}
		act := pt1.Add(pt2)
		if exp != act {
			t.Fatalf("expected %v + %v = %v but saw %v", pt1, pt2, exp, act)
		}
	}
}

// TestSub tests that two points can be subtracted.
func TestSub(t *testing.T) {
	rng := rand.New(rand.NewSource(22))
	for i := 0; i < 100; i++ {
		pt1 := Point{
			X: rng.Float64()*21 - 11,
			Y: rng.Float64()*21 - 11,
		}
		pt2 := Point{
			X: rng.Float64()*21 - 11,
			Y: rng.Float64()*21 - 11,
		}
		exp := Point{
			X: pt1.X - pt2.X,
			Y: pt1.Y - pt2.Y,
		}
		act := pt1.Sub(pt2)
		if exp != act {
			t.Fatalf("expected %v - %v = %v but saw %v", pt1, pt2, exp, act)
		}
	}
}

// TestMul tests that a point can be multiplied by a scalar.
func TestMul(t *testing.T) {
	rng := rand.New(rand.NewSource(33))
	for i := 0; i < 100; i++ {
		pt1 := Point{
			X: rng.Float64()*21 - 11,
			Y: rng.Float64()*21 - 11,
		}
		k := rng.Float64()*21 - 11
		exp := Point{
			X: pt1.X * k,
			Y: pt1.Y * k,
		}
		act := pt1.Mul(k)
		if exp != act {
			t.Fatalf("expected %v * %v = %v but saw %v", pt1, k, exp, act)
		}
	}
}

// TestDiv tests that a point can be divided by a scalar.
func TestDiv(t *testing.T) {
	rng := rand.New(rand.NewSource(44))
	for i := 0; i < 100; i++ {
		pt1 := Point{
			X: rng.Float64()*21 - 11,
			Y: rng.Float64()*21 - 11,
		}
		var k float64
		for k == 0.0 {
			k = rng.Float64()*21 - 11
		}
		exp := Point{
			X: pt1.X / k,
			Y: pt1.Y / k,
		}
		act := pt1.Div(k)
		if exp != act {
			t.Fatalf("expected %v / %v = %v but saw %v", pt1, k, exp, act)
		}
	}
}

// TestEq tests whether two Points can be compared for equality.
func TestEq(t *testing.T) {
	rng := rand.New(rand.NewSource(55))
	for tol := 0.0; tol < 1.0; tol += 0.1 {
		for i := 0; i < 100; i++ {
			pt1 := Point{
				X: rng.Float64()*21 - 11,
				Y: rng.Float64()*21 - 11,
			}
			pt2 := Point{
				X: pt1.X + tol*0.99,
				Y: pt1.Y - tol*0.99,
			}
			if !pt1.Eq(pt2, tol) {
				t.Fatalf("Points %v and %v should be equal within a tolerance of %v but aren't", pt1, pt2, tol)
			}
		}
	}
}

// TestFormat tests different ways of formatting a point.
func TestFormat(t *testing.T) {
	// Define a set of test cases in terms of a format string and its
	// expected output.
	pt := Point{X: 12.3, Y: -45.6}
	fmtExp := [][2]string{
		{"%v", "[12.3, -45.6]"},
		{"%#v", "xmorph.Point{X:12.3, Y:-45.6}"},
		{"%T", "xmorph.Point"},
		{"%f", "[12.300000, -45.600000]"},
		{"%.1f", "[12.3, -45.6]"},
		{"%6.1f", "[  12.3,  -45.6]"},
		{"%-10.2e", "[1.23e+01  , -4.56e+01 ]"},
		{"%07g", "[00012.3, -0045.6]"},
	}

	// Test each case in turn.
	for _, test := range fmtExp {
		format, exp := test[0], test[1]
		s := fmt.Sprintf(format, pt)
		if s != exp {
			t.Fatalf("Sprintf(%q, pt) produced %q, not %q", format, s, exp)
		}
	}
}
