// This file provides a mesh abstraction.

package morph

/*
#include <xmorph/mesh.h>
#include <xmorph/mesh_t.h>
#cgo LDFLAGS: -lmorph
*/
import "C"
import (
	"image"
	"runtime"
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
	mesh C.MeshT // Underlying mesh representation
}

// NewMesh creates a new mesh of given number of vertices.
func NewMesh(nx, ny int) *Mesh {
	m := &Mesh{}
	C.meshAlloc(&m.mesh, C.int(nx), C.int(ny))
	runtime.SetFinalizer(m, func(m *Mesh) {
		C.meshUnref(&m.mesh)
	})
	return m
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
