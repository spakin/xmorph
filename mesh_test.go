// The functions defined in this file ensure the morph package's test
// operations work as expected.

package morph

import (
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
	_, _, _, _ = m1, m2, m3, m4
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

// TestMeshFromImagePoints ensures we can create meshes from various-sized
// 2-D slices of image.Points.
func TestMeshFromImagePoints(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	s1 := random2DImagePoints(rng, 4, 4) // Minimal mesh size
	m1 := MeshFromImagePoints(s1)
	s2 := random2DImagePoints(rng, 101, 7) // Wide and short
	m2 := MeshFromImagePoints(s2)
	s3 := random2DImagePoints(rng, 7, 101) // Tall and narrow
	m3 := MeshFromImagePoints(s3)
	s4 := random2DImagePoints(rng, 3000, 3000) // Large
	m4 := MeshFromImagePoints(s4)
	_, _, _, _ = m1, m2, m3, m4
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

// TestMeshFromPoints ensures we can create meshes from various-sized
// 2-D slices of morph.Points.
func TestMeshFromPoints(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	s1 := random2DPoints(rng, 4, 4) // Minimal mesh size
	m1 := MeshFromPoints(s1)
	s2 := random2DPoints(rng, 101, 7) // Wide and short
	m2 := MeshFromPoints(s2)
	s3 := random2DPoints(rng, 7, 101) // Tall and narrow
	m3 := MeshFromPoints(s3)
	s4 := random2DPoints(rng, 3000, 3000) // Large
	m4 := MeshFromPoints(s4)
	_, _, _, _ = m1, m2, m3, m4
}
