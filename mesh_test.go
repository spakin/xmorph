// The functions defined in this file ensure the xmorph package's Mesh
// operations work as expected.

package xmorph

import (
	"bytes"
	"image"
	"math"
	"math/rand"
	"strings"
	"testing"
)

// TestNewEmptyMesh ensures that meshes of various sizes can be created and don't
// crash.
func TestNewEmptyMesh(t *testing.T) {
	m1 := NewEmptyMesh(4, 4)         // Minimal mesh size
	m2 := NewEmptyMesh(101, 7)       // Wide and short
	m3 := NewEmptyMesh(7, 101)       // Tall and narrow
	m4 := NewEmptyMesh(30000, 30000) // Very large
	m1.Free()
	m2.Free()
	m3.Free()
	m4.Free()
}

// random2DPoints allocates and populates a 2-D slice with random xmorph.Points.
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

// comparePointSlices tests that two 2-D slices of xmorph.Points are the same
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
// 2-D slices of xmorph.Points.
func TestMeshFromPoints(t *testing.T) {
	// Create a set of slices.
	rng := rand.New(rand.NewSource(11))
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

// compareImageePointSlices tests that two 2-D slices of xmorph.Points are the
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
	rng := rand.New(rand.NewSource(22))
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

// TestNewRegularMesh tests that we can create a regular mesh.
func TestNewRegularMesh(t *testing.T) {
	// Define a set of tests to perform.
	const wd, ht = 1024, 768
	nxNy := []image.Point{
		{9, 9}, {10, 10}, {12, 9}, {16, 18}, {102, 76},
	}
	close := func(a, b float64) bool {
		// Return true if two floating-point values are close enough to
		// be considered equal.
		return math.Abs(a-b) < 1e-10
	}

	// Ensure each mesh is indeed regular.
	for _, mSize := range nxNy {
		nx, ny := mSize.X, mSize.Y
		deltaX, deltaY := (wd-1)/float64(nx-1), (ht-1)/float64(ny-1)
		m := NewRegularMesh(nx, ny, wd, ht)
		sl := m.Points()
		for r := 0; r < ny-1; r++ {
			for c := 0; c < nx-1; c++ {
				// Test the x delta.
				pt := sl[r][c]
				dx := sl[r][c+1].Sub(pt).X
				if !close(dx, deltaX) {
					t.Fatalf("expected a delta of %.5g from (%d, %d) - (%d, %d) but saw %.5g",
						deltaX, c, r, c+1, r, dx)
				}

				// Test the y delta.
				dy := sl[r+1][c].Sub(pt).Y
				if !close(dy, deltaY) {
					t.Fatalf("expected a delta of %.5g from (%d, %d) - (%d, %d) but saw %.5g",
						deltaY, c, r, c, r+1, dy)
				}
			}
		}
		m.Free()
	}
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
	defer m.Free()

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
6 7
</orig>
<rect>
0 0 5 6
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
6 7
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

// TestReadMesh ensures we can read a mesh file into a Mesh.
func TestReadMesh(t *testing.T) {
	meshStr := `M2
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
6 7
</orig>
<rect>
0 0 5 6
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
6 7
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

	// Define a function that checks the mesh values.
	validate := func(m *Mesh) {
		sl := m.Points()
		for j := range sl {
			for i := range sl[j] {
				pt := Point{
					X: math.Round(float64(i)*12.5) / 10.0,
					Y: math.Round(float64(j)*17.5) / 10.0,
				}
				if sl[j][i] != pt {
					t.Fatalf("expected (%d, %d) = %v but observed %v", i, j, pt, sl[j][i])
				}
			}
		}
	}

	// Ensure we can read the mesh and receive the values we expect.
	m, err := ReadMesh(strings.NewReader(meshStr))
	if err != nil {
		t.Fatal(err)
	}
	validate(m)

	// Ensure we can read a mesh that lacks subimage data.
	lines := strings.Split(meshStr, "\n")
	m, err = ReadMesh(strings.NewReader(strings.Join(lines[:22], "\n")))
	if err != nil {
		t.Fatal(err)
	}
	validate(m)

	// Ensure that any other subset returns an error.
	for n := 21; n >= 0; n-- {
		_, err = ReadMesh(strings.NewReader(strings.Join(lines[:n], "\n")))
		if err == nil {
			t.Fatalf("truncating a mesh file to %d lines was supposed to return an error but didn't", n)
		}
	}
}

// TestMeshGetSet ensures we can get and set mesh points.
func TestMeshGetSet(t *testing.T) {
	// Set mesh points to arbitrary values.  We do so in column-major order
	// so the underlying vectors are not written in linear fashion.
	const (
		nx  = 11
		ny  = 10
		inc = 1.625
	)
	m := NewEmptyMesh(nx, ny)
	idx := 1.0
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			var pt Point
			pt.X = idx
			idx += inc
			pt.Y = idx
			idx += inc
			m.Set(i, j, pt)
		}
	}

	// Read back all values previously written.
	idx = 1.0
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			var pt Point
			pt.X = idx
			idx += inc
			pt.Y = idx
			idx += inc
			mpt := m.Get(i, j)
			if mpt != pt {
				t.Fatalf("wrote %v to (%d, %d) but read back %v", pt, i, j, mpt)
			}
		}
	}
}

// TestMeshGetSetImage ensures we can get and set mesh points as image.Point
// values.
func TestMeshGetSetImage(t *testing.T) {
	// Set mesh points to arbitrary values.  We do so in column-major order
	// so the underlying vectors are not written in linear fashion.
	const (
		nx  = 10
		ny  = 11
		inc = 3
	)
	m := NewEmptyMesh(nx, ny)
	idx := 1
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			var pt image.Point
			pt.X = idx
			idx += inc
			pt.Y = idx
			idx += inc
			m.SetImagePoint(i, j, pt)
		}
	}

	// Read back all values previously written.
	idx = 1.0
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			var pt image.Point
			pt.X = idx
			idx += inc
			pt.Y = idx
			idx += inc
			mpt := m.GetImagePoint(i, j)
			if mpt != pt {
				t.Fatalf("wrote %v to (%d, %d) but read back %v", pt, i, j, mpt)
			}
		}
	}
}

// TestFunctionalize ensures we can functionalize a mesh.
func TestFunctionalize(t *testing.T) {
	// Test a slice containing completely random coordinates.
	// (random2DPoints is too structures for this test.)
	const mwd, mht = 8, 6   // Mesh width and height
	const wd, ht = 800, 600 // Image width and height
	rng := rand.New(rand.NewSource(44))
	sl := make([][]Point, mht)
	for r := range sl {
		row := make([]Point, mwd)
		for c := range row {
			row[c] = Point{
				X: rng.Float64() * wd,
				Y: rng.Float64() * ht,
			}
		}
		sl[r] = row
	}
	m := MeshFromPoints(sl)
	const expected = 12 // Empirically determined
	actual := m.Functionalize(wd, ht)
	if actual != expected {
		t.Fatalf("expected functionalization to fix %d points, but it fixed %d", expected, actual)
	}
}

// validateMeshDimens checks that a Mesh has the expected dimensions and that
// these are consistent across the C and Go versions.
func validateMeshDimens(t *testing.T, m *Mesh, wd, ht int) {
	if int(m.mesh.nx) != m.NX || int(m.mesh.ny) != m.NY {
		t.Fatalf("inconsistent mesh dimensions (%d, %d) in Go vs. (%d, %d) from C", m.NX, m.NY, int(m.mesh.nx), int(m.mesh.ny))
	}
	if m.NX != wd || m.NY != ht {
		t.Fatalf("invalid mesh dimensions: expected (%d, %d) but saw (%d, %d)", wd, ht, m.NX, m.NY)
	}
}

// TestScale ensures we can scale mesh coordinates.
func TestScale(t *testing.T) {
	// Create a mesh.
	m := NewEmptyMesh(5, 5)
	fs := []float64{0.0, 0.25, 0.5, 0.75, 1.0}
	for r, y := range fs {
		for c, x := range fs {
			m.Set(c, r, Point{X: x, Y: y})
		}
	}

	// Scale it to a 100x50 image.
	m.Scale(100, 50)
	for r, y := range fs {
		for c, x := range fs {
			exp := Point{x * 100.0, y * 50.0}
			if exp.X == 100.0 {
				exp.X = 99.0 // Clamp to image bounds.
			}
			if exp.Y == 50.0 {
				exp.Y = 49.0 // Clamp to image bounds.
			}
			act := m.Get(c, r)
			if act != exp {
				t.Fatalf("expected %v at (%d, %d) but saw %v", exp, c, r, act)
			}
		}
	}
}

// TestAddLine ensures we can add a line to a mesh.
func TestAddLine(t *testing.T) {
	// Create a mesh.
	m := NewEmptyMesh(5, 5)
	fs := []float64{0.0, 10.0, 20.0, 50.0, 100.0}
	for r, y := range fs {
		for c, x := range fs {
			m.Set(c, r, Point{X: x, Y: y})
		}
	}

	// Add a vertical line to the left.
	err := m.AddLine(0, 0.25, Vertical)
	if err != nil {
		t.Fatal(err)
	}

	// Add a vertical line to the right.
	err = m.AddLine(4, 0.5, Vertical)
	if err != nil {
		t.Fatal(err)
	}

	// Add a horizontal line to the bottom.
	err = m.AddLine(3, 0.33, Horizontal)
	if err != nil {
		t.Fatal(err)
	}

	// Add a horizontal line to the middle.
	err = m.AddLine(1, 0.82, Horizontal)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the results.
	if m.NX != 7 || m.NY != 7 {
		t.Fatalf("expected the augmented mesh to be 7x7 but saw %dx%d", m.NX, m.NY)
	}
	exp := [][]Point{
		{{0, 0}, {2.5, 0}, {10, 0}, {20, 0}, {50, 0}, {75, 0}, {100, 0}},
		{{0, 10}, {2.5, 10}, {10, 10}, {20, 10}, {50, 10}, {75, 10}, {100, 10}},
		{{0, 18.2}, {2.5, 18.2}, {10, 18.2}, {20, 18.2}, {50, 18.2}, {75, 18.2}, {100, 18.2}},
		{{0, 20}, {2.5, 20}, {10, 20}, {20, 20}, {50, 20}, {75, 20}, {100, 20}},
		{{0, 50}, {2.5, 50}, {10, 50}, {20, 50}, {50, 50}, {75, 50}, {100, 50}},
		{{0, 66.5}, {2.5, 66.5}, {10, 66.5}, {20, 66.5}, {50, 66.5}, {75, 66.5}, {100, 66.5}},
		{{0, 100}, {2.5, 100}, {10, 100}, {20, 100}, {50, 100}, {75, 100}, {100, 100}},
	}
	act := m.Points()
	for r, row := range exp {
		for c := range row {
			if act[r][c] != exp[r][c] {
				t.Fatalf("expected %v but saw %v at (%d, %d)", exp[r][c], act[r][c], c, r)
			}
		}
	}
}

// TestDeleteLine ensures we can delete a line to a mesh.
func TestDeleteLine(t *testing.T) {
	// Create a mesh.
	m := NewEmptyMesh(6, 6)
	fs := []float64{0.0, 10.0, 20.0, 50.0, 85.0, 100.0}
	for r, y := range fs {
		for c, x := range fs {
			m.Set(c, r, Point{X: x, Y: y})
		}
	}

	// Remove the leftmost column.
	err := m.DeleteLine(0, Vertical)
	if err != nil {
		t.Fatal(err)
	}

	// Remove the bottommost row.
	err = m.DeleteLine(5, Horizontal)
	if err != nil {
		t.Fatal(err)
	}

	// Remove the central column.
	err = m.DeleteLine(2, Vertical)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the results.
	if m.NX != 4 || m.NY != 5 {
		t.Fatalf("expected the augmented mesh to be 7x7 but saw %dx%d", m.NX, m.NY)
	}
	exp := [][]Point{
		{{10, 0}, {20, 0}, {85, 0}, {100, 0}},
		{{10, 10}, {20, 10}, {85, 10}, {100, 10}},
		{{10, 20}, {20, 20}, {85, 20}, {100, 20}},
		{{10, 50}, {20, 50}, {85, 50}, {100, 50}},
		{{10, 85}, {20, 85}, {85, 85}, {100, 85}},
	}
	act := m.Points()
	for r, row := range exp {
		for c := range row {
			if act[r][c] != exp[r][c] {
				t.Fatalf("expected %v but saw %v at (%d, %d)", exp[r][c], act[r][c], c, r)
			}
		}
	}
}

// TestCopy ensures we can deep-copy a mesh.
func TestCopy(t *testing.T) {
	// Ensure that no data changes during a copy.
	rng := rand.New(rand.NewSource(55))
	const wd, ht = 25, 25
	sl := random2DPoints(rng, wd, ht)
	m1 := MeshFromPoints(sl)
	m2 := m1.Copy()
	for j := 0; j < ht; j++ {
		for i := 0; i < wd; i++ {
			pt1 := m1.Get(i, j)
			pt2 := m2.Get(i, j)
			if pt1 != pt2 {
				t.Fatalf("mismatch at (%d, %d): expected %v but saw %v", i, j, pt1, pt2)
			}
		}
	}

	// Ensure that the NX and NY fields were copied, too.  We also confirm
	// that they were assigned correctly to begin with because we don't
	// currently have a separate test for that.
	validateMeshDimens(t, m1, wd, ht)
	validateMeshDimens(t, m2, wd, ht)

	// Ensure that the copy was deep.  If we change an element in the
	// source, it should not change in the target.
	cx, cy := wd/2, ht/2
	vOld := m1.Get(cx, cy)
	vNew := vOld.Mul(2.0)
	m1.Set(cx, cy, vNew)
	for j := 0; j < ht; j++ {
		for i := 0; i < wd; i++ {
			pt1 := m1.Get(i, j)
			pt2 := m2.Get(i, j)
			if pt1 != pt2 && (i != cx || j != cy) {
				t.Fatalf("mismatch at (%d, %d): expected %v but saw %v", i, j, pt1, pt2)
			}
			if pt1 == pt2 && i == cx && j == cy {
				t.Fatalf("unexpected match at (%d, %d): expected %v and %v but saw only %v",
					i, j, vOld, vNew, pt1)
			}
		}
	}
}

// TestInterpolate ensures we can interpolate two meshes.
func TestInterpolate(t *testing.T) {
	// Create two meshes: One regular and one with all internal points
	// scaled downwards.
	const wd, ht = 10, 8
	const scale = 0.25
	const interp = 0.6
	sl1 := make([][]Point, ht)
	for r := range sl1 {
		row := make([]Point, wd)
		for c := range row {
			row[c] = Point{
				X: float64(c * 100),
				Y: float64(r * 100),
			}
		}
		sl1[r] = row
	}
	m1 := MeshFromPoints(sl1)
	sl2 := make([][]Point, ht)
	for r := 0; r < ht; r++ {
		sl2[r] = make([]Point, wd)
		if r == 0 || r == ht-1 {
			copy(sl2[r], sl1[r])
			continue
		}
		sl2[r][0] = sl1[r][0]
		for c := 1; c < wd-1; c++ {
			sl2[r][c] = sl1[r][c].Mul(scale)
		}
		sl2[r][wd-1] = sl1[r][wd-1]
	}
	m2 := MeshFromPoints(sl2)

	// Interpolate the two meshes.
	mi, err := InterpolateMeshes(m1, m2, interp)
	if err != nil {
		t.Fatal(err)
	}
	sli := mi.Points()

	// Ensure all meshes have the same dimensions.
	validateMeshDimens(t, m1, wd, ht)
	validateMeshDimens(t, m2, wd, ht)
	validateMeshDimens(t, mi, wd, ht)

	// Check the results.
	for r := 0; r < ht; r++ {
		for c := 0; c < wd; c++ {
			var expected Point
			if r == 0 || r == ht-1 || c == 0 || c == wd-1 {
				// Edges should be unmodified.
				expected = sl1[r][c]
			} else {
				// Edges should be interpolated.
				pt1 := sl1[r][c]
				pt2 := sl2[r][c]
				expected = pt1.Mul(1.0 - interp).Add(pt2.Mul(interp))
			}
			actual := sli[r][c]
			if expected != actual {
				t.Fatalf("failed interpolation at (%d, %d): expected %v but saw %v", c, r, expected, actual)
			}
		}
	}
}
