// The functions defined in this file ensure the xmorph package's warp
// operations work as expected.

package xmorph

import (
	"crypto/sha256"
	"encoding/base64"
	"image"
	"os"
	"strings"
	"testing"

	"image/color"
	"image/png"
	_ "image/png"
)

// gopherString is a base64-encoded 128x128 PNG image of the Go gopher.
const gopherString = `
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
o6Oqvrz8fzomEF6Q0gdGAAAAAElFTkSuQmCC
`

// gopherMeshIn is a 5x5 input mesh to apply to the image of the Go gopher.
var gopherMeshIn = MeshFromPoints([][]Point{
	{{0, 0}, {32, 0}, {64, 0}, {96, 0}, {128, 0}},
	{{0, 32}, {32, 32}, {64, 32}, {96, 32}, {128, 32}},
	{{0, 64}, {32, 64}, {64, 64}, {96, 64}, {128, 64}},
	{{0, 96}, {32, 96}, {64, 96}, {96, 96}, {128, 96}},
	{{0, 128}, {32, 128}, {64, 128}, {96, 128}, {128, 128}},
})

// gopherMeshOut is a 5x5 output mesh to apply to the image of the Go gopher.
var gopherMeshOut = MeshFromPoints([][]Point{
	{{0, 0}, {32, 0}, {64, 0}, {96, 0}, {128, 0}},
	{{0, 32}, {26, 23}, {64, 32}, {112, 10}, {128, 32}},
	{{0, 64}, {32, 64}, {64, 64}, {96, 64}, {128, 64}},
	{{0, 96}, {17, 114}, {64, 96}, {116, 112}, {128, 96}},
	{{0, 128}, {32, 128}, {64, 128}, {96, 128}, {128, 128}},
})

// gopherImage is a 128x128 image of the Go gopher, expanded from gopherString.
var gopherImage image.Image

// init expands gopherString to gopherImage.
func init() {
	r := strings.NewReader(gopherString)
	dec := base64.NewDecoder(base64.StdEncoding, r)
	var err error
	gopherImage, _, err = image.Decode(dec)
	if err != nil {
		panic(err)
	}
}

// copyImage copies a given image's data, converting the color model as it
// goes.
func copyImage(cm color.Model, set func(x, y int, c color.Color), img image.Image) {
	bnds := gopherImage.Bounds()
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			c := img.At(x, y)
			set(x, y, cm.Convert(c))
		}
	}
}

// imageHash computes a SHA256 hash of an image.
func imageHash(t *testing.T, img image.Image) []byte {
	hash := sha256.New()
	err := png.Encode(hash, img)
	if err != nil {
		t.Fatal(err)
	}
	return hash.Sum(nil)
}

// compareHashes returns an error if an expected hash and an actual hash are
// not equal.
func compareHashes(t *testing.T, h1, h2 []byte) {
	if len(h1) != len(h2) {
		t.Fatalf("hashes are of different lengths: %v vs. %v", h1, h2)
	}
	for i := range h1 {
		if h1[i] != h2[i] {
			t.Fatalf("hash mismatch: expected %#v but saw %#v", h1, h2)
		}
	}
}

// writePNG writes an image to a PNG file.  This function can be used to
// validate visually a generated image.
func writePNG(fn string, img image.Image) {
	f, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

// TestWarpNRGBA tests that an NRGBA image can be warped according to a source
// and destination mesh.
func TestWarpNRGBA(t *testing.T) {
	// Warp the image.
	img := image.NewNRGBA(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x8e, 0x5e, 0x2c, 0x5d, 0x74, 0xdb, 0xcd, 0x37, 0xd0,
		0x8c, 0xdc, 0x33, 0x8d, 0xe2, 0x7d, 0x27, 0x9a, 0xd9, 0x6e,
		0xb6, 0xfc, 0x95, 0xeb, 0x99, 0xc4, 0xc2, 0xcb, 0x64, 0x2e,
		0x20, 0xad, 0xf2}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarpAlpha tests that an Alpha image can be warped according to a source
// and destination mesh.
func TestWarpAlpha(t *testing.T) {
	// Warp the image.
	img := image.NewAlpha(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x41, 0x3f, 0x79, 0x1a, 0xfc, 0x1e, 0x2a, 0x4a, 0x8c,
		0x2, 0xe5, 0x25, 0x6b, 0x71, 0x67, 0xa0, 0xe7, 0xc5, 0x2a,
		0x26, 0x6b, 0x62, 0x16, 0xb9, 0xcd, 0x6d, 0xeb, 0xe5, 0x18,
		0x40, 0x5f, 0x64}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarpCMYK tests that a CMYK image can be warped according to a source
// and destination mesh.
func TestWarpCMYK(t *testing.T) {
	// Warp the image.
	img := image.NewCMYK(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x19, 0x27, 0x24, 0x5e, 0xbe, 0xac, 0xfb, 0x72, 0xf9,
		0xe7, 0x1e, 0xd1, 0x72, 0x28, 0x44, 0x87, 0x66, 0x1, 0xb5,
		0x45, 0x37, 0x4c, 0xfe, 0x73, 0xb5, 0xeb, 0xf9, 0xb4, 0x81,
		0xdc, 0x87, 0x9f}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarpGray tests that a Gray image can be warped according to a source
// and destination mesh.
func TestWarpGray(t *testing.T) {
	// Warp the image.
	img := image.NewGray(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0xf5, 0x6b, 0x26, 0xf3, 0xb4, 0x2e, 0xc9, 0xff, 0xf6,
		0x82, 0xd8, 0xa7, 0xa2, 0xc8, 0xae, 0x9a, 0x19, 0x50, 0x70,
		0xd1, 0x81, 0xc6, 0x8e, 0x11, 0xe4, 0xb3, 0xc5, 0x53, 0x3f,
		0x3f, 0xb1, 0x4e}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarpRGBA tests that an RGBA image can be warped according to a source
// and destination mesh.
func TestWarpRGBA(t *testing.T) {
	// Warp the image.
	img := image.NewRGBA(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x1a, 0x17, 0xca, 0x47, 0xcd, 0xf9, 0xf2, 0xea, 0xc5,
		0x2a, 0x7a, 0xfd, 0x31, 0x7b, 0xf4, 0x1e, 0x22, 0xfa, 0xab,
		0x46, 0xab, 0x95, 0x21, 0xd, 0x82, 0xb9, 0x41, 0xe, 0xd7,
		0x6a, 0x2f, 0x9a}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarpGray16 tests that a Gray16 image can be warped according to a source
// and destination mesh.
func TestWarpGray16(t *testing.T) {
	// Warp the image.
	img := image.NewGray16(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x7d, 0x39, 0x1b, 0x4b, 0xfe, 0x31, 0x17, 0x46, 0x4,
		0xfe, 0x41, 0xf4, 0x72, 0x8c, 0xa3, 0x26, 0x60, 0xae, 0x6b,
		0x43, 0xe, 0x7f, 0x0, 0x69, 0x3, 0xc8, 0x6d, 0xd4, 0xe1,
		0x54, 0x90, 0x9}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarpRGBA64 tests that an RGBA64 image can be warped according to a
// source and destination mesh.
func TestWarpRGBA64(t *testing.T) {
	// Warp the image.
	img := image.NewRGBA64(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0xe7, 0x99, 0x77, 0xe3, 0xfc, 0x21, 0x68, 0xe6, 0xc2,
		0x59, 0xe6, 0xe3, 0xa7, 0xd6, 0xf5, 0xd0, 0xfb, 0x7a, 0xdd,
		0x38, 0x4b, 0xa7, 0x2c, 0x2b, 0xfd, 0xbc, 0xee, 0xa5, 0x10,
		0x29, 0x5a, 0x24}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarp0NRGBA tests that warping an NRGBA image 0% of the way from a source
// to a destination mesh does not noticeably change the image.
func TestWarp0NRGBA(t *testing.T) {
	// Warp the image.
	img := image.NewNRGBA(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 0.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x9, 0x84, 0x12, 0x10, 0x93, 0x72, 0x96, 0xa2, 0x16,
		0x7b, 0x95, 0x2a, 0x78, 0xbc, 0xc8, 0x0, 0x23, 0xb, 0x54,
		0xf5, 0x73, 0xd8, 0x82, 0x5c, 0x22, 0x7d, 0xb9, 0xe8, 0xff,
		0x6d, 0x9d, 0xfc}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarp25NRGBA tests that warping an NRGBA image 25% of the way from a
// source to a destination mesh produces the expected output.
func TestWarp25NRGBA(t *testing.T) {
	// Warp the image.
	img := image.NewNRGBA(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 0.25)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x6f, 0x6b, 0x7a, 0xd6, 0xf7, 0x58, 0xb8, 0x51, 0xd1,
		0x49, 0x3d, 0xbf, 0x83, 0x1d, 0xb0, 0xcd, 0xdd, 0x95, 0xe9,
		0x27, 0xca, 0xf1, 0x6f, 0xc5, 0x22, 0x69, 0x82, 0x55, 0xa8,
		0x23, 0xee, 0x35}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarp66NRGBA tests that warping an NRGBA image 66% of the way from a
// source to a destination mesh produces the expected output.
func TestWarp66NRGBA(t *testing.T) {
	// Warp the image.
	img := image.NewNRGBA(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 0.66)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x25, 0xe7, 0xfd, 0xfe, 0xd3, 0x20, 0x4e, 0x5a, 0xa3,
		0xab, 0x4b, 0xee, 0xab, 0xf5, 0xde, 0x37, 0xef, 0x41, 0x88,
		0x9d, 0xaf, 0xd4, 0x7f, 0xd6, 0x7e, 0x62, 0x3b, 0x4f, 0x1a,
		0xef, 0x8b, 0xe2}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)
}

// TestWarpNRGBANN tests that an NRGBA image can be warped according to a
// source and destination mesh and using nearest-neighbor interpolation instead
// of Lanczos-based antialiasing.
func TestWarpNRGBANN(t *testing.T) {
	// Select the antialiasing kernel.
	Antialiasing = NearestNeighbor

	// Warp the image.
	img := image.NewNRGBA(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0xef, 0x55, 0xfb, 0xa6, 0x26, 0x66, 0x8b, 0xed, 0xbe,
		0x4c, 0x1d, 0xbc, 0x77, 0xbd, 0x8e, 0xe1, 0xf8, 0xab, 0x99,
		0x14, 0x15, 0x93, 0x10, 0x90, 0x2c, 0xa9, 0xf6, 0x3c, 0xdc,
		0xf5, 0xff, 0xc}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)

	// Restore the default Lanczos antialiasing.
	Antialiasing = Lanczos
}

// TestWarpNRGBABilinear tests that an NRGBA image can be warped according to a
// source and destination mesh and using bilinear interpolation instead of
// Lanczos-based antialiasing.
func TestWarpNRGBABilinear(t *testing.T) {
	// Select the antialiasing kernel.
	Antialiasing = Bilinear

	// Warp the image.
	img := image.NewNRGBA(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x22, 0x18, 0x81, 0x63, 0x6c, 0x17, 0x7d, 0x63, 0x5d,
		0x64, 0xd8, 0xa8, 0xd5, 0xf2, 0x4, 0x7, 0x8a, 0x8d, 0xee,
		0x63, 0x71, 0x7, 0x41, 0xf, 0xff, 0xd5, 0xef, 0xdd, 0x46,
		0xa9, 0x4b, 0x89}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)

	// Restore the default Lanczos antialiasing.
	Antialiasing = Lanczos
}

// TestWarpNRGBALanczos4 tests that an NRGBA image can be warped according to a
// source and destination mesh and using the higher-quality Lanczos-based
// antialiasing.
func TestWarpNRGBALanczos4(t *testing.T) {
	// Select the antialiasing kernel.
	Antialiasing = Lanczos4

	// Warp the image.
	img := image.NewNRGBA(gopherImage.Bounds())
	copyImage(img.ColorModel(), img.Set, gopherImage)
	warp, err := Warp(img, gopherMeshIn, gopherMeshOut, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the image's hash value to an expected value.
	exp := []byte{0x29, 0x8b, 0xa, 0xe8, 0x30, 0xee, 0xa5, 0x5e, 0xba,
		0xb9, 0xb3, 0x96, 0x49, 0x10, 0xdd, 0xa3, 0x30, 0x0, 0xc3,
		0xb6, 0x33, 0xb3, 0xdb, 0xb9, 0xc4, 0xc, 0x47, 0x43, 0x51,
		0x31, 0xac, 0xc9}
	hash := imageHash(t, warp)
	compareHashes(t, exp, hash)

	// Restore the default Lanczos antialiasing.
	Antialiasing = Lanczos
}
