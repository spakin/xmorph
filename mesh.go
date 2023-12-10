// This file provides a mesh abstraction.

package xmorph

/*
#include <stdlib.h>
#include <xmorph/mesh.h>
#include <xmorph/mesh_t.h>
#cgo LDFLAGS: -lmorph
*/
import "C"
import (
	"bufio"
	"fmt"
	"image"
	"io"
	"math"
	"strings"
	"unsafe"
)

// A Mesh represents a 2-D mesh.
type Mesh struct {
	NX   int      // Number of mesh points in the x direction
	NY   int      // Number of mesh points in the y direction
	mesh *C.MeshT // Underlying mesh representation
}

// NewEmptyMesh creates a new, empty mesh of a given number of vertices.
func NewEmptyMesh(nx, ny int) *Mesh {
	return &Mesh{
		NX:   nx,
		NY:   ny,
		mesh: C.meshNew(C.int(nx), C.int(ny)),
	}
}

// Free deallocates the mesh memory managed by libmorph.
func (m *Mesh) Free() {
	C.meshUnref(m.mesh)
	m.mesh = nil
}

// MeshFromPoints creates a new mesh from a 2-D slice of xmorph.Points.
func MeshFromPoints(sl [][]Point) *Mesh {
	// Sanity check the mesh lest libmorph write something itself to
	// standard error.
	if len(sl) < 4 || len(sl[0]) < 4 {
		panic("slice passed to MeshFromPoints must be at least 4x4")
	}
	nx, ny := len(sl[0]), len(sl)
	for _, row := range sl {
		if len(row) != nx {
			panic("all rows in the MeshFromPoints slice must be the same length")
		}
	}

	// Create an empty mesh.
	m := NewEmptyMesh(nx, ny)

	// Populate the mesh element by element.
	np := nx * ny
	xp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.x))[:np:np]
	yp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.y))[:np:np]
	i := 0
	for _, row := range sl {
		for _, pt := range row {
			xp[i] = C.double(pt.X)
			yp[i] = C.double(pt.Y)
			i++
		}
	}
	return m
}

// ReadMesh reads a morph/xmorph/gtkmorph mesh file and returns a Mesh object.
func ReadMesh(r io.Reader) (*Mesh, error) {
	// An empty error from a scan implies EOF.  We want to report
	// it as such.
	scanner := bufio.NewScanner(r)
	getErr := func() error {
		err := scanner.Err()
		if err == nil {
			return io.EOF
		}
		return err
	}

	// Parse the file header.
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read the mesh header (%w)", getErr())
	}
	if scanner.Text() != "M2" {
		return nil, fmt.Errorf("invalid mesh header (should be \"M2\")")
	}

	// Read the mesh dimensions.
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read the mesh dimensions (%w)", getErr())
	}
	var nx, ny int
	ln := scanner.Text()
	ntoks, err := fmt.Sscan(ln, &nx, &ny)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q as mesh dimensions (%w)", ln, err)
	}
	if ntoks != 2 {
		return nil, fmt.Errorf("failed to parse %q as mesh dimensions", ln)
	}
	if nx < 4 || ny < 4 {
		return nil, fmt.Errorf("mesh must be at least 4x4 (read %dx%d)", nx, ny)
	}

	// Parse each of the remaining lines into mesh coordinates and a label.
	sl := make([][]Point, ny)
	for j := range sl {
		sl[j] = make([]Point, nx)
		for i := range sl[j] {
			if !scanner.Scan() {
				return nil, fmt.Errorf("failed to read %d mesh coordinates (%w)", nx*ny, getErr())
			}
			ln = scanner.Text()
			var x, y int
			var label int
			ntoks, err = fmt.Sscan(ln, &x, &y, &label)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %q as {x, y, label} (%w)", ln, err)
			}
			if ntoks != 3 {
				return nil, fmt.Errorf("failed to parse %q as {x, y, label}", ln)
			}
			sl[j][i] = Point{
				X: float64(x) / 10.0,
				Y: float64(y) / 10.0,
			}
		}
	}

	// Create and return a Mesh object.
	return MeshFromPoints(sl), nil
}

// MeshFromImagePoints creates a new mesh from a 2-D slice of image.Points.
func MeshFromImagePoints(sl [][]image.Point) *Mesh {
	// Sanity check the mesh lest libmorph doesn't write something itself
	// to standard error.
	if len(sl) < 4 || len(sl[0]) < 4 {
		panic("slice passed to MeshFromImagePoints must be at least 4x4")
	}
	nx, ny := len(sl[0]), len(sl)
	for _, row := range sl {
		if len(row) != nx {
			panic("all rows in the MeshFromImagePoints slice must be the same length")
		}
	}

	// Create an empty mesh.
	m := NewEmptyMesh(nx, ny)

	// Populate the mesh element by element.
	np := nx * ny
	xp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.x))[:np:np]
	yp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.y))[:np:np]
	i := 0
	for _, row := range sl {
		for _, pt := range row {
			xp[i] = C.double(pt.X)
			yp[i] = C.double(pt.Y)
			i++
		}
	}
	return m
}

// NewRegularMesh creates a new, regular mesh of a given number of vertices.
func NewRegularMesh(nx, ny, wd, ht int) *Mesh {
	wd1, ht1 := wd-1, ht-1
	nx1, ny1 := float64(nx-1), float64(ny-1)
	sl := make([][]Point, ny)
	for r := range sl {
		sl[r] = make([]Point, nx)
		y := float64(r*ht1) / ny1
		for c := range sl[r] {
			x := float64(c*wd1) / nx1
			sl[r][c] = Point{X: x, Y: y}
		}
	}
	return MeshFromPoints(sl)
}

// Points converts a mesh to a 2-D slice of xmorph.Points.
func (m *Mesh) Points() [][]Point {
	// Represent the MeshT's flat lists of x and y values as Go slices.
	nx, ny := int(m.mesh.nx), int(m.mesh.ny)
	np := nx * ny
	xp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.x))[:np:np]
	yp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.y))[:np:np]

	// Reshape the flat lists as a Go slice of slices.
	sl := make([][]Point, ny)
	idx := 0
	for j := range sl {
		sl[j] = make([]Point, nx)
		for i := range sl[j] {
			sl[j][i].X = float64(xp[idx])
			sl[j][i].Y = float64(yp[idx])
			idx++
		}
	}
	return sl
}

// ImagePoints converts a mesh to a 2-D slice of image.Points.
func (m *Mesh) ImagePoints() [][]image.Point {
	// Represent the MeshT's flat lists of x and y values as Go slices.
	nx, ny := int(m.mesh.nx), int(m.mesh.ny)
	np := nx * ny
	xp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.x))[:np:np]
	yp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.y))[:np:np]

	// Reshape the flat lists as a Go slice of slices.
	sl := make([][]image.Point, ny)
	idx := 0
	for j := range sl {
		sl[j] = make([]image.Point, nx)
		for i := range sl[j] {
			sl[j][i].X = int(xp[idx])
			sl[j][i].Y = int(yp[idx])
			idx++
		}
	}
	return sl
}

// meshRanges returns the x and y ranges of the mesh data.
func (m *Mesh) meshRanges() (Point, Point) {
	// Represent the MeshT's flat lists of x and y values as Go slices.
	nx, ny := int(m.mesh.nx), int(m.mesh.ny)
	np := nx * ny
	xp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.x))[:np:np]
	yp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.y))[:np:np]

	// Find the coordinates of the upper left and lower right
	// corners of the mesh.
	ul := Point{X: float64(xp[0]), Y: float64(yp[0])}
	lr := ul
	for i := 0; i < np; i++ {
		x, y := float64(xp[i]), float64(yp[i])
		ul.X = math.Min(ul.X, x)
		ul.Y = math.Min(ul.Y, y)
		lr.X = math.Max(lr.X, x)
		lr.Y = math.Max(lr.Y, y)
	}
	return ul, lr
}

// Write outputs a mesh that's compatible with morph, xmorph, and gtkmorph.
func (m *Mesh) Write(w io.Writer) error {
	// Write the two header lines.
	var err error
	if _, err = fmt.Fprintln(w, "M2"); err != nil {
		return err
	}
	nx, ny := int(m.mesh.nx), int(m.mesh.ny)
	if _, err = fmt.Fprintln(w, nx, ny); err != nil {
		return err
	}

	// Write all of the data.
	np := nx * ny
	xp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.x))[:np:np]
	yp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.y))[:np:np]
	lp := (*[1 << 30]C.int)(unsafe.Pointer(m.mesh.label))[:np:np]
	for i := range xp {
		x := math.Round(float64(xp[i]) * 10.0)
		y := math.Round(float64(yp[i]) * 10.0)
		_, err = fmt.Fprintf(w, "%.0f %.0f %d\n", x, y, lp[i])
		if err != nil {
			return err
		}
	}

	// Write the subimage information.
	ul, lr := m.meshRanges()
	dx, dy := lr.X-ul.X, lr.Y-ul.Y
	eye1 := Point{X: ul.X + dx/3, Y: ul.Y + dy/3}
	eye2 := Point{X: ul.X + 2*dx/3, Y: ul.Y + dy/3}
	eye3 := Point{X: ul.X + dx/2, Y: ul.Y + 2*dy/3}
	_, err = fmt.Fprintf(w, `<SIS>
<orig>
%.0f %.0f
</orig>
<rect>
%.0f %.0f %.0f %.0f
</rect>
<eye>
%.6f %.6f
</eye>
<eye>
%.6f %.6f
</eye>
<eye>
%.6f %.6f
</eye>
</SIS>
<resulting image size>
%.0f %.0f
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
`,
		math.Ceil(dx+1), math.Ceil(dy+1),
		math.Floor(ul.X), math.Floor(ul.Y),
		math.Ceil(lr.X), math.Ceil(lr.Y),
		eye1.X, eye1.Y,
		eye2.X, eye2.Y,
		eye3.X, eye3.Y,
		math.Ceil(dx+1), math.Ceil(dy+1))
	if err != nil {
		return err
	}
	return nil
}

// checkMeshCoord panics if a given coordinate lies out of range.
func (m *Mesh) checkMeshCoord(x, y int) {
	cx, cy := C.int(x), C.int(y)
	if cx < 0 || cy < 0 || C.long(cx) >= m.mesh.nx || C.long(cy) >= m.mesh.ny {
		panic(fmt.Sprintf("point (%d, %d) lies of out bounds of the mesh", x, y))
	}
}

// Get returns the xmorph.Point at (x, y).
func (m *Mesh) Get(x, y int) Point {
	m.checkMeshCoord(x, y)
	cx, cy := C.int(x), C.int(y)
	var pt Point
	pt.X = float64(C.meshGetx(m.mesh, cx, cy))
	pt.Y = float64(C.meshGety(m.mesh, cx, cy))
	return pt
}

// GetImagePoint returns the image.Point at (x, y).
func (m *Mesh) GetImagePoint(x, y int) image.Point {
	pt := m.Get(x, y)
	return image.Point{
		X: int(pt.X),
		Y: int(pt.Y),
	}
}

// Set assigns the xmorph.Point at (x, y).
func (m *Mesh) Set(x, y int, pt Point) {
	m.checkMeshCoord(x, y)
	cx, cy := C.int(x), C.int(y)
	C.meshSetNoundo(m.mesh, cx, cy, C.double(pt.X), C.double(pt.Y))
}

// SetImagePoint assigns the image.Point at (x, y).
func (m *Mesh) SetImagePoint(x, y int, pt image.Point) {
	m.checkMeshCoord(x, y)
	cx, cy := C.int(x), C.int(y)
	C.meshSetNoundo(m.mesh, cx, cy, C.double(pt.X), C.double(pt.Y))
}

// Functionalize fixes problems with the mesh, making it both functional and
// bounded.  It takes as input the width and height of the image to which the
// mesh corresponds and returns the number of changes made.
func (m *Mesh) Functionalize(w, h int) int {
	nc := C.meshFunctionalize(m.mesh, C.int(w), C.int(h))
	return int(nc)
}

// Scale scales mesh coordinates to fit a given image width and height.
func (m *Mesh) Scale(w, h int) {
	C.meshScale(m.mesh, C.int(w), C.int(h))
}

// A Direction can be either horizontal or vertical.
type Direction int

// These are the acceptable values for a direction.
const (
	Vertical   Direction = 1
	Horizontal           = 2
)

// AddLine adds a row or column to the mesh, fraction f of the way
// from index i to index i + 1.
func (m *Mesh) AddLine(i int, f float64, d Direction) error {
	// Sanity-check our arguments so libmorph doesn't write its own error
	// message to stderr.
	switch d {
	case Vertical:
		if i < 0 || i >= m.NX-1 {
			return fmt.Errorf("index %d lies outside the range [0, %d]", i, m.NX-2)
		}
	case Horizontal:
		if i < 0 || i >= m.NY-1 {
			return fmt.Errorf("index %d lies outside the range [0, %d]", i, m.NY-2)
		}
	default:
		return fmt.Errorf("unexpected direction %d", d)
	}
	if f < 0.0 || f > 1.0 {
		return fmt.Errorf("line-adding fraction must lie in [0.0, 1.0]")
	}

	// Add the line.
	r := C.meshLineAdd(m.mesh, C.int(i), C.double(f), C.int(d))
	if r != 0 {
		return fmt.Errorf("AddLine failed to add a line (id = %d)", r)
	}
	m.NX = int(m.mesh.nx)
	m.NY = int(m.mesh.ny)
	return nil
}

// DeleteLine deletes a row or column from the mesh.
func (m *Mesh) DeleteLine(i int, d Direction) error {
	// Sanity-check our arguments so libmorph doesn't write its own error
	// message to stderr.
	switch d {
	case Vertical:
		if i < 0 || i >= m.NX {
			return fmt.Errorf("index %d lies outside the range [0, %d]", i, m.NX-1)
		}
	case Horizontal:
		if i < 0 || i >= m.NY {
			return fmt.Errorf("index %d lies outside the range [0, %d]", i, m.NY-1)
		}
	default:
		return fmt.Errorf("unexpected direction %d", d)
	}

	// Delete the line.
	r := C.meshLineDelete(m.mesh, C.int(i), C.int(d))
	if r != 0 {
		return fmt.Errorf("AddLine failed to add a line (id = %d)", r)
	}
	m.NX = int(m.mesh.nx)
	m.NY = int(m.mesh.ny)
	return nil
}

// Copy deep-copies a mesh.
func (m *Mesh) Copy() *Mesh {
	mc := NewEmptyMesh(int(m.mesh.nx), int(m.mesh.ny))
	C.meshCopy(mc.mesh, m.mesh)
	return mc
}

// Format returns a textual representation of a Mesh's coordinates, using the
// flags it receives to format each Point.
func (m *Mesh) Format(st fmt.State, verb rune) {
	frags := make([]string, m.NY)
	pts := m.Points()
	for r, row := range pts {
		rFrags := make([]string, m.NX)
		for c, pt := range row {
			rFrags[c] = pt.formatString(st, verb)
		}
		rStr := strings.Join(rFrags, ", ")
		frags[r] = fmt.Sprintf("[%s]", rStr)
	}
	fmt.Fprintf(st, "[%s]", strings.Join(frags, ", "))
}

// InterpolateMeshes interpolates two meshes to produce a new mesh that lies a
// given fraction from the first mesh's points to the second mesh's points.  It
// returns an error code if the meshes are incompatible.
func InterpolateMeshes(m1, m2 *Mesh, t float64) (*Mesh, error) {
	if t < 0.0 || t > 1.0 {
		return nil, fmt.Errorf("interpolation fraction %.5g does not lie in the range [0.0, 1.0]", t)
	}
	if C.meshCompatibilityCheck(m1.mesh, m2.mesh) != 0 {
		return nil, fmt.Errorf("incompatible meshes passed to InterpolateMeshes")
	}
	m := m1.Copy()
	C.meshInterpolate(m.mesh, m1.mesh, m2.mesh, C.double(t))
	return m, nil
}
