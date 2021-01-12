// This file provides a mesh abstraction.

package morph

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
	"unsafe"
)

// A Point is a float64-valued (x, y) coordinate.  The axes increase right and
// down.
type Point struct {
	X float64
	Y float64
}

// A Mesh represents a 2-D mesh.
type Mesh struct {
	mesh *C.MeshT // Underlying mesh representation
}

// NewMesh creates a new mesh of a given number of vertices.
func NewMesh(nx, ny int) *Mesh {
	m := &Mesh{}
	m.mesh = C.meshNew(C.int(nx), C.int(ny))
	return m
}

// Free deallocates the mesh memory managed by libmorph.
func (m *Mesh) Free() {
	C.meshUnref(m.mesh)
	m.mesh = nil
}

// MeshFromPoints creates a new mesh from a 2-D slice of morph.Points.
func MeshFromPoints(sl [][]Point) *Mesh {
	// Sanity check the mesh lest libmorph doesn't write something itself
	// to standard error.
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
	m := NewMesh(nx, ny)

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
	m := NewMesh(nx, ny)

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

// Points converts a mesh to a 2-D slice of morph.Points.
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

// Get returns the morph.Point at (x, y).
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

// Set assigns the morph.Point at (x, y).
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

// Copy deep-copies a mesh.
func (m *Mesh) Copy() *Mesh {
	mc := NewMesh(int(m.mesh.nx), int(m.mesh.ny))
	C.meshCopy(mc.mesh, m.mesh)
	return mc
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
