// This file provides various usage examples of the xmorph package.

package xmorph_test

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"strings"

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

// Warp an image (hard-wired in this example) from a regular source mesh 75% of
// the way to a randomly perturbed destination mesh.
func ExampleWarp() {
	// Create a small image.
	imgStr := `
/9j/4AAQSkZJRgABAQIA7ADsAAD/2wBDABALDA4MChAODQ4SERATGCgaGBYWGDEjJR0oOjM9PDkz
ODdASFxOQERXRTc4UG1RV19iZ2hnPk1xeXBkeFxlZ2P/wgALCABAADIBAREA/8QAGgAAAgMBAQAA
AAAAAAAAAAAAAAECAwQFBv/aAAgBAQAAAAH0FFUI9EzS5+mG+XPdUndpMANdB5KUPolOQH0CrE1o
0syl1h//xAAfEAACAQQDAQEAAAAAAAAAAAABAgMAEBESEyEyIEH/2gAIAQEAAQUCppAKGXaQ6PEX
5LSPiiuIuXicybigcinjIYMRWo26tG3eaIyL46Hq0o7v+2l8XX1aTxdEPwyjbioKFv8A/8QAIhAA
AQMDBAMBAAAAAAAAAAAAAQAQEQIgIRIxUXETMkFh/9oACAEBAAY/AlH1ZKFNIq5lZeBuj441fqAP
twHltS3WVhoaLJQ7ebB3eHNkm3dYf//EACIQAAIBBAICAwEAAAAAAAAAAAERABAhMVFBYSCxcZGh
8f/aAAgBAQABPyGGUvFqIuVAiMSLBtLAbZt3Vb+KEECxbYw5zPIIVKQ5mLjMV7U0MzzuAkCAgI5l
HMQNENUU3GRE3HZcy4KOakgHA4lhekSoEBvY1UH4PdQZ/I9+AMfdRfhBYnwDBYBl1bIonA7g2+om
OX3X/9oACAEBAAAAEGlAyCWAixz/xAAlEAEAAQMDBAEFAAAAAAAAAAABEQAhYRAxUUFxkcGhIIGx
0fD/2gAIAQEAAT8QpchxzYqeGJIUdimUXSEg3GevWhQouBbLzI9dRVlbCmH5pINTezJTW/8AcAGN
6UATQBVcbFCoSAyPDRjEQmOMabpNrFl03OZtXD+CBipBGZMxvPehEtmZMpo7TLGAmHqe/NDghhzR
9APDSICBZONUEHVDSuSWLMb290AABBY0JdlxcsSPw/GiCIkjQDYD7VC//I1zAQ8PoyQXxf1rAOIf
CNb7ae9s0hYQjos4aKKQELzq4ENhf05pcfafmavndLuv/9k=
`
	r := strings.NewReader(imgStr)
	dec := base64.NewDecoder(base64.StdEncoding, r)
	img, _, err := image.Decode(dec)
	if err != nil {
		panic(err)
	}

	// Define an regular, initial mesh and a randomly warped mesh.
	bnds := img.Bounds()
	wd := bnds.Max.X - bnds.Min.X
	ht := bnds.Max.Y - bnds.Min.Y
	const nx, ny = 5, 12
	mReg := xmorph.NewRegularMesh(nx, ny, wd, ht)
	mWarp := mReg.Copy()
	dx, dy := float64(wd)/nx/2.0, float64(ht)/ny/2.0
	rng := rand.New(rand.NewSource(12345))
	for r := 1; r < ny-1; r++ {
		for c := 1; c < nx-1; c++ {
			pt := mWarp.Get(c, r)
			pt.X += rng.Float64()*dx*2.0 - dx
			pt.Y += rng.Float64()*dy*2.0 - dy
			mWarp.Set(c, r, pt)
		}
	}

	// Warp the image 75% of the way from the initial mesh to the warped mesh.
	wImg, err := xmorph.Warp(img, mReg, mWarp, 0.75)
	if err != nil {
		panic(err)
	}
	enc := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	opt := jpeg.Options{Quality: 50}
	err = jpeg.Encode(enc, wImg, &opt)
	if err != nil {
		panic(err)
	}
	// Output:
	// /9j/2wCEABALDA4MChAODQ4SERATGCgaGBYWGDEjJR0oOjM9PDkzODdASFxOQERXRTc4UG1RV19iZ2hnPk1xeXBkeFxlZ2MBERISGBUYLxoaL2NCOEJjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY//AAAsIAEAAMgEBEQD/xADSAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgsQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+v/aAAgBAQAAPwDv6hnuYbcFppVX0BPWsmW4lvJQoLKrsAFU9RS6j/oEkCWcMwaXd+9DEouMH5s+pP6UugST3KfaJhIm5fmjfsf8f8a2qKzdSv8AylaO3OZcgM2Puf8A16dbQrLpzLE3kzPGVaReWViOufXvWP8A2gdK1aOynWS5uTCHMp4B9e3XirMl012VctlG5UL0Uf41Ja3bWpAALRE5Kjt6kf4VsqwZQykMp6EHrTq5q/sGS7SRi4RZnlXYeGz6+9TQzywyExSBA3JBXPNQyq086zStvfjeuPlf29uKlfYXYxRCNDyFX+vvSZxzVvS59jSwHOxcMmBxg9cfQ/zrSE0RH3x+dMubcXEe1iVPUEdjWNIjRSujjDA/5NJRStG6Irsvyv0IOQaIjsuoJNzAbwpwcDDYFbwRAAABxTj0rE1Aq2oyAA7ljUNx65/PtUNIQCCCTg9xSKCsYj3MUHRSeBQQGZFPUuoH13CujorP1W2UxG6BCyxrj2dfT/D3rOPBopKltE8y+hQ/wnf+X/663aKztWSR442AJjRsso9ex/Cs3r0pR0o749at2FtK8n2gqAqrhQ3f3HpWkqSBQDLkgelS0VlXdqpuG8p9rvjEZ4BPfH86QaZOTy8aDuRlj+HSrUOnQxYLAysOhft9BVzFFf/Z
}

// Morph a red square 25% of the way to a magenta triangle.
func ExampleMorph() {
	// Define a 32x32 image of a red square.
	sqrStr := `
iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAIAAAD8GO2jAAAAPUlEQVRIx2P8//8/Ay0BEwONwagF
oxZQDlhwSXxiYibJIL5/f0fjYNSCUQtGLRi1YOhawDjaLhq1YNQCBgB6cwk7256ScAAAAABJRU5E
rkJggg==
`
	r := strings.NewReader(sqrStr)
	dec := base64.NewDecoder(base64.StdEncoding, r)
	sqrImg, _, err := image.Decode(dec)
	if err != nil {
		panic(err)
	}

	// Define a 32x32 image of a magenta triangle.
	triStr := `
iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAIAAAD8GO2jAAABAElEQVRIx2P8//8/Ay0BEwONwagF
oxbQ24Kfz798PP2ahhY86Hp/t5xEL/wnGnx/8vkg05/9DP/fH31JvC4SfHC/9cO/f8wMDAz3ahiJ
18VIZFn07f6n08o8cAfpHX4pZCNOzTi43/IZ2bsPapmo6YOvtz6cUef/z4ASMkR6giiH3G/6hmY6
8Z4grOjzlXdvlkphin86IPru0EsqWHC/4SfObFHHRKkFny68ebdWEqfsQdF3B19QZMG9mr8E/FfH
Qr4FH0+//rCVQDr5fEjk7f4XZFpwr4aoPHi/hpUcC94dfPFxlxgxFnw5Jvx23wuSLSAYuER6gnG0
XTRqwagFDABVVMHneY8tAgAAAABJRU5ErkJggg==
`
	r = strings.NewReader(triStr)
	dec = base64.NewDecoder(base64.StdEncoding, r)
	triImg, _, err := image.Decode(dec)
	if err != nil {
		panic(err)
	}

	// Define a mesh that wraps the square.
	sqrMesh := xmorph.MeshFromImagePoints([][]image.Point{
		{{0, 0}, {10, 0}, {21, 0}, {31, 0}},
		{{0, 8}, {8, 8}, {24, 8}, {31, 8}},
		{{0, 24}, {8, 24}, {24, 24}, {31, 24}},
		{{0, 31}, {10, 31}, {21, 31}, {31, 31}},
	})

	// Define a mesh that wraps the triangle.
	triMesh := xmorph.MeshFromImagePoints([][]image.Point{
		{{0, 0}, {10, 0}, {21, 0}, {31, 0}},
		{{0, 8}, {15, 8}, {16, 8}, {31, 8}},
		{{0, 24}, {8, 24}, {24, 24}, {31, 24}},
		{{0, 31}, {10, 31}, {21, 31}, {31, 31}},
	})

	// Morph the image 25% of the way from the square to the triangle.
	mImg, err := xmorph.Morph(sqrImg, triImg, sqrMesh, triMesh, 0.25)
	if err != nil {
		panic(err)
	}
	b64Enc := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	pEnc := png.Encoder{CompressionLevel: png.BestCompression}
	pEnc.Encode(b64Enc, mImg)
	if err != nil {
		panic(err)
	}
	// Output:
	// iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAC6ElEQVR42uxXQU8bPRB943X08VGVYy/9Ib33R/TU38q96rEHqBpAEUlJlQpQQhEk9kxl7zqZ9dpJb1xYaXdjz9jz/PxmvLEiInjBy4SHBiHq0n2l3/mVbIfm0z6kBxERaYfU3jore+6Xj98HUvtRjjafJAdSA3mobxswA2tqA/WqD21BrS8PXGLK5Mhqe76P4h6lap58QaVtNaWJSnuvbXmAksgOAR1kQWmCHMi/Bs7tegE5w3Yf2jSYJxMsP30WbFrBg6J4EaYhAxhrCA0BBIhnyJrx5vQU9uQ/dL4D6gcASvR0feKmC3ocb8D3Rsz/JvImToSfBWQE5tiLPSaQJfJPjM0dg778prcf34PsUMg6li0F7wtOyI+n4m6PyJER84hkIA5WD5gV4FZAiBBuLyLubA7/4R01J1ZK6VethKXS8fR1DAFj1DKOdDeAhDspmbrfIwCb83OsF3Vxp5i2hm7XFnLfrsFoQsAYpH3FR/JPfbHtAjtn47YX3EJVcw62oLQNCZRsAD+5ClsJSSYFAru+ZI/s8MV34addptfqiCkeEDtGQrEm/FrFdUQ6FN36LZ0tAZH5gsxRmItpX8U0B4oG8eIOTM9oiMLEKn5v1RKhdu0m4HfPkOUDJfp1Kd5bCQdH8XwmglEQ+3a1A6xd8M4uLdAGZjmVVgPDIjY4jGpa4MsJBBtqdsH1W1S7pxHGBn58SSlE7awxtdNr2399AaJR0gBQ3gZRWSEc0rWxoJsfGt9g9YPTUKdjB0b8xQziXHAk6addDCzdFkChMnGgA1/NpJ8sw8Wa2qfU9lrcEHFkQGhH+bYOaDqSLTyMsSK3c6plWY8BvfJcA266FHZ+y4AOqoSgxNeKkdmT/7mqpl9PhJr2vjPD3/8J9SoefSnPoouWYDoZu3Y8L2HBD4/QWVASoq0ViYhQCOspw8Agq5OiWaeuOLYQYz2SRoD1zAPccVf4Nhx8lL7Y/4JXAK8AXgG85PU3AAD//8XhrcD86lqbAAAAAElFTkSuQmCC
}
