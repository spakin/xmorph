// The functions defined in this file ensure the morph package's test
// operations work as expected.

package morph

import (
	"bytes"
	"image"
	"math/rand"
	"testing"
)

// TestNewMesh ensures that meshes of various sizes can be created and don't
// crash.
func TestNewMesh(t *testing.T) {
	m1 := NewMesh(4, 4)         // Minimal mesh size
	m2 := NewMesh(101, 7)       // Wide and short
	m3 := NewMesh(7, 101)       // Tall and narrow
	m4 := NewMesh(30000, 30000) // Very large
	m1.Free()
	m2.Free()
	m3.Free()
	m4.Free()
}

// random2DPoints allocates and populates a 2-D slice with random morph.Points.
func random2DPoints(rng *rand.Rand, nx, ny int) [][]Point {
	s := make([][]Point, ny)
	for j := range s {
		s[j] = make([]Point, nx)
		var pt Point
		for i := range s[j] {
			// Use monotically increasing x values.
			if i == 0 {
				pt.X = rng.Float64() * 10
			} else {
				pt.X = s[j][i-1].X + rng.Float64()*10 + 1
			}

			// Use monotically increasing y values.
			if j == 0 {
				pt.Y = rng.Float64() * 10
			} else {
				pt.Y = s[j-1][i].Y + rng.Float64()*10 + 1
			}

			// Assign the point.
			s[j][i] = pt
		}
	}
	return s
}

// comparePointSlices tests that two 2-D slices of morph.Points are the same
// size and contain the same values.  It aborts if not.
func comparePointSlices(t *testing.T, s1, s2 [][]Point) {
	if len(s1) != len(s2) {
		t.Fatalf("mismatched row counts (%d vs. %d)", len(s1), len(s2))
	}
	for j, row := range s1 {
		if len(s1[j]) != len(s2[j]) {
			t.Fatalf("mismatched column counts in row %d (%d vs. %d)", j, len(s1[j]), len(s2[j]))
		}
		for i := range row {
			if s1[j][i] != s2[j][i] {
				t.Fatalf("mismatched values at s[%d][%d]: %v vs. %v", j, i, s1[j][i], s2[j][i])
			}
		}
	}
}

// TestMeshFromPoints ensures we can create meshes from various-sized
// 2-D slices of morph.Points.
func TestMeshFromPoints(t *testing.T) {
	// Create a set of slices.
	rng := rand.New(rand.NewSource(22))
	i1 := random2DPoints(rng, 4, 4)       // Minimal mesh size
	i2 := random2DPoints(rng, 101, 7)     // Wide and short
	i3 := random2DPoints(rng, 7, 101)     // Tall and narrow
	i4 := random2DPoints(rng, 3000, 3000) // Large

	// Convert slices to meshes.
	m1 := MeshFromPoints(i1)
	m2 := MeshFromPoints(i2)
	m3 := MeshFromPoints(i3)
	m4 := MeshFromPoints(i4)

	// Convert meshes back to slices.
	o1 := m1.Points()
	o2 := m2.Points()
	o3 := m3.Points()
	o4 := m4.Points()

	// Ensure the before and after slices match.
	comparePointSlices(t, i1, o1)
	comparePointSlices(t, i2, o2)
	comparePointSlices(t, i3, o3)
	comparePointSlices(t, i4, o4)

	// Deallocate libmorph resources.
	m1.Free()
	m2.Free()
	m3.Free()
	m4.Free()
}

// random2DImagePoints allocates and populates a 2-D slice with random
// image.Points.
func random2DImagePoints(rng *rand.Rand, nx, ny int) [][]image.Point {
	s := make([][]image.Point, ny)
	for j := range s {
		s[j] = make([]image.Point, nx)
		var pt image.Point
		for i := range s[j] {
			// Use monotically increasing x values.
			if i == 0 {
				pt.X = rng.Intn(10)
			} else {
				pt.X = s[j][i-1].X + rng.Intn(10) + 1
			}

			// Use monotically increasing y values.
			if j == 0 {
				pt.Y = rng.Intn(10)
			} else {
				pt.Y = s[j-1][i].Y + rng.Intn(10) + 1
			}

			// Assign the point.
			s[j][i] = pt
		}
	}
	return s
}

// compareImageePointSlices tests that two 2-D slices of morph.Points are the
// same size and contain the same values.  It aborts if not.
func compareImagePointSlices(t *testing.T, s1, s2 [][]image.Point) {
	if len(s1) != len(s2) {
		t.Fatalf("mismatched row counts (%d vs. %d)", len(s1), len(s2))
	}
	for j, row := range s1 {
		if len(s1[j]) != len(s2[j]) {
			t.Fatalf("mismatched column counts in row %d (%d vs. %d)", j, len(s1[j]), len(s2[j]))
		}
		for i := range row {
			if s1[j][i] != s2[j][i] {
				t.Fatalf("mismatched values at s[%d][%d]: %v vs. %v", j, i, s1[j][i], s2[j][i])
			}
		}
	}
}

// TestMeshFromImagePoints ensures we can create meshes from various-sized
// 2-D slices of image.Points.
func TestMeshFromImagePoints(t *testing.T) {
	// Create a set of slices.
	rng := rand.New(rand.NewSource(11))
	i1 := random2DImagePoints(rng, 4, 4)       // Minimal mesh size
	i2 := random2DImagePoints(rng, 101, 7)     // Wide and short
	i3 := random2DImagePoints(rng, 7, 101)     // Tall and narrow
	i4 := random2DImagePoints(rng, 3000, 3000) // Large

	// Convert slices to meshes.
	m1 := MeshFromImagePoints(i1)
	m2 := MeshFromImagePoints(i2)
	m3 := MeshFromImagePoints(i3)
	m4 := MeshFromImagePoints(i4)

	// Convert meshes back to slices.
	o1 := m1.ImagePoints()
	o2 := m2.ImagePoints()
	o3 := m3.ImagePoints()
	o4 := m4.ImagePoints()

	// Ensure the before and after slices match.
	compareImagePointSlices(t, i1, o1)
	compareImagePointSlices(t, i2, o2)
	compareImagePointSlices(t, i3, o3)
	compareImagePointSlices(t, i4, o4)

	// Deallocate libmorph resources.
	m1.Free()
	m2.Free()
	m3.Free()
	m4.Free()
}

// TestWrite tests that we can write a mesh to an io.Writer.
func TestWrite(t *testing.T) {
	// Create a Mesh with known contents.
	nx, ny := 5, 4
	s := make([][]Point, ny)
	for j := range s {
		s[j] = make([]Point, nx)
		for i := range s[j] {
			s[j][i].X = float64(i) * 1.25
			s[j][i].Y = float64(j) * 1.75
		}
	}
	m := MeshFromPoints(s)

	// Write a mesh file to a bytes.Buffer.
	var buf bytes.Buffer
	m.Write(&buf)

	// Ensure the result is as expected.
	expected := `M2
5 4
0 0 0
13 0 0
25 0 0
38 0 0
50 0 0
0 18 0
13 18 0
25 18 0
38 18 0
50 18 0
0 35 0
13 35 0
25 35 0
38 35 0
50 35 0
0 53 0
13 53 0
25 53 0
38 53 0
50 53 0
<SIS>
<orig>
5 5
</orig>
<rect>
0 0 5 5
</rect>
<eye>
1.666667 1.750000
</eye>
<eye>
3.333333 1.750000
</eye>
<eye>
2.500000 3.500000
</eye>
</SIS>
<resulting image size>
5 5
</resulting image size>
<features>
<name>
feature 0
</name>
<name>
feature 1
</name>
<name>
feature 2
</name>
</features>
`
	actual := buf.String()
	if actual != expected {
		t.Logf(actual)
		t.Fatalf("unexpected output (expected %d bytes; observed %d bytes)", len(expected), len(actual))
	}
}
