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
func MeshFromPoints(s [][]Point) *Mesh {
	// Sanity check the mesh lest libmorph doesn't write something itself
	// to standard error.
	if len(s) < 4 || len(s[0]) < 4 {
		panic("slice passed to MeshFromPoints must be at least 4x4")
	}
	nx, ny := len(s[0]), len(s)
	for _, row := range s {
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
	for _, row := range s {
		for _, pt := range row {
			xp[i] = C.double(pt.X)
			yp[i] = C.double(pt.Y)
			i++
		}
	}
	return m
}

// MeshFromImagePoints creates a new mesh from a 2-D slice of image.Points.
func MeshFromImagePoints(s [][]image.Point) *Mesh {
	// Sanity check the mesh lest libmorph doesn't write something itself
	// to standard error.
	if len(s) < 4 || len(s[0]) < 4 {
		panic("slice passed to MeshFromImagePoints must be at least 4x4")
	}
	nx, ny := len(s[0]), len(s)
	for _, row := range s {
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
	for _, row := range s {
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
	s := make([][]Point, ny)
	idx := 0
	for j := range s {
		s[j] = make([]Point, nx)
		for i := range s[j] {
			s[j][i].X = float64(xp[idx])
			s[j][i].Y = float64(yp[idx])
			idx++
		}
	}
	return s
}

// ImagePoints converts a mesh to a 2-D slice of image.Points.
func (m *Mesh) ImagePoints() [][]image.Point {
	// Represent the MeshT's flat lists of x and y values as Go slices.
	nx, ny := int(m.mesh.nx), int(m.mesh.ny)
	np := nx * ny
	xp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.x))[:np:np]
	yp := (*[1 << 30]C.double)(unsafe.Pointer(m.mesh.y))[:np:np]

	// Reshape the flat lists as a Go slice of slices.
	s := make([][]image.Point, ny)
	idx := 0
	for j := range s {
		s[j] = make([]image.Point, nx)
		for i := range s[j] {
			s[j][i].X = int(xp[idx])
			s[j][i].Y = int(yp[idx])
			idx++
		}
	}
	return s
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
