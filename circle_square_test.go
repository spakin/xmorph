// This file presents a complete morphing example.

package xmorph_test

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"

	"github.com/spakin/xmorph"
)

const nFrames = 29 // Number of animation frames to generate (odd avoids repetition in the middle)

// Abs returns the absolute value of an integer.
func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

// DrawCircle creates an image of a colored circle.
func DrawCircle(wd, ht, cx, cy, r int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, wd, ht))
	r2 := r * r
	for y := 0; y < ht; y++ {
		for x := 0; x < wd; x++ {
			dist := (x-cx)*(x-cx) + (y-cy)*(y-cy)
			if dist > r2 {
				// Outside the circle: set to black.
				img.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 255})
			} else {
				// Within the circle: set to a color between
				// green and white.
				rb := 255 - uint8(255*dist/r2)
				img.SetNRGBA(x, y, color.NRGBA{rb, 255, rb, 255})
			}
		}
	}
	return img
}

// DrawSquare creates an image of a colored square.
func DrawSquare(wd, ht, cx, cy, e int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, wd, ht))
	e2 := e / 2
	for y := 0; y < ht; y++ {
		for x := 0; x < wd; x++ {
			if x >= cx-e2 && x < cx+e2 && y >= cy-e2 && y < cy+e2 {
				// Within the square: set to a color between
				// yellow and white.
				dist := Abs(cx-x) + Abs(cy-y)
				b := 255 - uint8(255*dist/e)
				img.SetNRGBA(x, y, color.NRGBA{255, 255, b, 255})
			} else {
				// Outside the square: set to black.
				img.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 255})
			}
		}
	}
	return img
}

// CreateCircleMesh wraps a mesh around a circle.
func CreateCircleMesh(wd, ht, cx, cy, r int) *xmorph.Mesh {
	// Start with an inscribed square, then expand the north, south, east
	// and west extrema to better envelop the circle.
	e := float64(r) / math.Sqrt2 // Projection of a 45 degree angle
	m := CreateSquareMesh(wd, ht, cx, cy, int(2*e))
	m.SetImagePoint(2, 1, image.Point{X: cx, Y: cy - r}) // North
	m.SetImagePoint(2, 3, image.Point{X: cx, Y: cy + r}) // South
	m.SetImagePoint(3, 2, image.Point{X: cx + r, Y: cy}) // East
	m.SetImagePoint(1, 2, image.Point{X: cx - r, Y: cy}) // West
	return m
}

// CreateSquareMesh wraps a mesh around a square.
func CreateSquareMesh(wd, ht, cx, cy, e int) *xmorph.Mesh {
	// Determine where all x and y coordinates should lie.
	e2 := e / 2
	xs := []int{0, cx - e2, cx, cx + e2 - 1, wd - 1}
	ys := []int{0, cy - e2, cy, cy + e2 - 1, ht - 1}

	// Produce a grid as the cross product of all x and all y coordinates.
	gr := make([][]image.Point, 5)
	for r, y := range ys {
		gr[r] = make([]image.Point, 5)
		for c, x := range xs {
			gr[r][c] = image.Point{X: x, Y: y}
		}
	}
	return xmorph.MeshFromImagePoints(gr)
}

// PrepareCircle returns an image of a circle and a corresponding mesh.
func PrepareCircle(wd, ht, cx, cy, r int) (image.Image, *xmorph.Mesh) {
	img := DrawCircle(wd, ht, cx, cy, r)
	mesh := CreateCircleMesh(wd, ht, cx, cy, r)
	return img, mesh
}

// PrepareSquare returns an image of a square and a corresponding mesh.
func PrepareSquare(wd, ht, cx, cy, r int) (image.Image, *xmorph.Mesh) {
	img := DrawSquare(wd, ht, cx, cy, r)
	mesh := CreateSquareMesh(wd, ht, cx, cy, r)
	return img, mesh
}

// GreenYellowPalette creates a color palette containing black, white, and a
// range of saturations of green and yellow.
func GreenYellowPalette() color.Palette {
	pal := make([]color.Color, 256)
	for i := 0; i < 128; i++ {
		c := uint8(255 * i / 127)
		pal[i] = color.NRGBA{255, 255, c, 255}
		pal[i+127] = color.NRGBA{c, 255, c, 255}
	}
	pal[255] = color.NRGBA{0, 0, 0, 255} // Replace a duplicate white with black.
	return pal
}

// MakePaletted converts an arbitrarily colored image to a paletted image.
func MakePaletted(img image.Image, pal color.Palette) *image.Paletted {
	bnds := img.Bounds()
	pImg := image.NewPaletted(bnds, pal)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			pImg.Set(x, y, img.At(x, y))
		}
	}
	return pImg
}

// This is a complete example of gradually morphing one image to another.
// Regrettably, the code is rather large because it is fully self-contained:
// all images and meshes are generated internally rather than read from files.
//
// The code morphs a large circle colored with a green and white gradient
// pattern and positioned in the upper-left quadrant of the image to a large
// square colored with a yellow and white gradient and positioned in the
// lower-right quadrant of the image.  The code saves the morph sequence to an
// animated GIF image called circle-square.gif.
func Example_animatedGIF() {
	// Create the source and destination image and mesh.
	cImg, cMesh := PrepareCircle(256, 256, 96, 96, 64)
	sImg, sMesh := PrepareSquare(256, 256, 160, 160, 128)

	// Create a sequence of frames in the forward direction and copy these
	// to the backward direction.
	frames := make([]*image.Paletted, nFrames)
	pal := GreenYellowPalette()
	nf2 := (nFrames + 1) / 2
	for i := 0; i < nf2; i++ {
		t := float64(i) / float64(nf2-1)
		img, err := xmorph.Morph(cImg, sImg, cMesh, sMesh, t)
		if err != nil {
			panic(err)
		}
		pImg := MakePaletted(img, pal)
		frames[i] = pImg
		frames[nFrames-1-i] = pImg
	}

	// Show all frames for 100 milliseconds except the middle frame, which
	// we show for 2 seconds and the (identical) first and last frames,
	// which we show for 1 second apiece.
	delay := make([]int, nFrames)
	for i := range delay {
		delay[i] = 10
	}
	delay[0] = 100
	delay[nFrames/2] = 200
	delay[nFrames-1] = 100

	// Define a disposal for each frame.
	disp := make([]byte, nFrames)
	for i := range disp {
		disp[i] = gif.DisposalNone
	}

	// Create an animated GIF.
	desc := gif.GIF{
		Image:     frames,
		Delay:     delay,
		LoopCount: 0,
		Disposal:  disp,
		Config: image.Config{
			Width:      256,
			Height:     256,
			ColorModel: pal,
		},
	}
	w, err := os.Create("circle-square.gif")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	err = gif.EncodeAll(w, &desc)
	if err != nil {
		panic(err)
	}
}
