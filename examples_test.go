// This file provides various usage examples of the xmorph package.

package xmorph_test

import (
	"fmt"
	"image"

	"github.com/spakin/xmorph"
)

// Define meshes m1 and m3 then interpolate the points 50% of the way from m1
// to m3 to produce mesh m2.
func ExampleInterpolateMeshes() {
	const nx, ny = 4, 4
	m1 := xmorph.NewRegularMesh(nx, ny, 100, 100)
	m3 := xmorph.NewRegularMesh(nx, ny, 100, 100)
	for r := 1; r < ny-1; r++ {
		for c := 1; c < nx-1; c++ {
			pt := m3.Get(c, r)
			pt.X /= 2.0
			pt.Y /= 2.0
			m3.Set(c, r, pt)
		}
	}
	m2, err := xmorph.InterpolateMeshes(m1, m3, 0.5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("m1 = %v\n", m1)
	fmt.Printf("m2 = %v\n", m2)
	fmt.Printf("m3 = %v\n", m3)
	// Output:
	// m1 = [[[0, 0], [33, 0], [66, 0], [99, 0]], [[0, 33], [33, 33], [66, 33], [99, 33]], [[0, 66], [33, 66], [66, 66], [99, 66]], [[0, 99], [33, 99], [66, 99], [99, 99]]]
	// m2 = [[[0, 0], [33, 0], [66, 0], [99, 0]], [[0, 33], [24.75, 24.75], [49.5, 24.75], [99, 33]], [[0, 66], [24.75, 49.5], [49.5, 49.5], [99, 66]], [[0, 99], [33, 99], [66, 99], [99, 99]]]
	// m3 = [[[0, 0], [33, 0], [66, 0], [99, 0]], [[0, 33], [16.5, 16.5], [33, 16.5], [99, 33]], [[0, 66], [16.5, 33], [33, 33], [99, 66]], [[0, 99], [33, 99], [66, 99], [99, 99]]]
}

func ExampleMeshFromPoints() {
	pts := make([][]xmorph.Point, 4)
	for r := range pts {
		pts[r] = make([]xmorph.Point, 4)
		for c := range pts[r] {
			pts[r][c] = xmorph.Point{
				X: float64(c) * 100.0,
				Y: float64(r) * 75.0,
			}
		}
	}
	m := xmorph.MeshFromPoints(pts)
	fmt.Println(m)
	// Output:
	// [[[0, 0], [100, 0], [200, 0], [300, 0]], [[0, 75], [100, 75], [200, 75], [300, 75]], [[0, 150], [100, 150], [200, 150], [300, 150]], [[0, 225], [100, 225], [200, 225], [300, 225]]]
}

func ExampleMeshFromImagePoints() {
	pts := make([][]image.Point, 4)
	for r := range pts {
		pts[r] = make([]image.Point, 4)
		for c := range pts[r] {
			pts[r][c] = image.Point{
				X: c * 100,
				Y: r * 75,
			}
		}
	}
	m := xmorph.MeshFromImagePoints(pts)
	fmt.Println(m)
	// Output:
	// [[[0, 0], [100, 0], [200, 0], [300, 0]], [[0, 75], [100, 75], [200, 75], [300, 75]], [[0, 150], [100, 150], [200, 150], [300, 150]], [[0, 225], [100, 225], [200, 225], [300, 225]]]
}
