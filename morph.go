// This file provides functions for morphing images.

package morph

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"reflect"
)

// avgU8 returns the weighted average of two uint8 values.
func avgU8(a, b uint8, t float64) uint8 {
	fa := float64(a) * (1.0 - t)
	fb := float64(b) * t
	return uint8(math.Round(fa + fb))
}

// morphNRGBA morphs an NRGBA image.
func morphNRGBA(sImg, dImg *image.NRGBA, sMesh, dMesh *Mesh, t float64) (*image.NRGBA, error) {
	// Separately warp the source and destination images.
	sw, err := Warp(sImg, sMesh, dMesh, t)
	if err != nil {
		return nil, err
	}
	dw, err := Warp(dImg, sMesh, dMesh, 1.0-t)
	if err != nil {
		return nil, err
	}
	sWarp, dWarp := sw.(*image.NRGBA), dw.(*image.NRGBA)

	// Perform a weighted average of the source and destination images'
	// colors to produce a final image.
	bnds := sWarp.Bounds()
	img := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			cs := sWarp.NRGBAAt(x, y)
			cd := dWarp.NRGBAAt(x, y)
			c := color.NRGBA{
				R: avgU8(cs.R, cd.R, t),
				G: avgU8(cs.G, cd.G, t),
				B: avgU8(cs.B, cd.B, t),
				A: avgU8(cs.A, cd.A, t),
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img, nil
}


// morphCMYK morphs a CMYK image.
func morphCMYK(sImg, dImg *image.CMYK, sMesh, dMesh *Mesh, t float64) (*image.CMYK, error) {
	// Separately warp the source and destination images.
	sw, err := Warp(sImg, sMesh, dMesh, t)
	if err != nil {
		return nil, err
	}
	dw, err := Warp(dImg, sMesh, dMesh, 1.0-t)
	if err != nil {
		return nil, err
	}
	sWarp, dWarp := sw.(*image.CMYK), dw.(*image.CMYK)

	// Perform a weighted average of the source and destination images'
	// colors to produce a final image.
	bnds := sWarp.Bounds()
	img := image.NewCMYK(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			cs := sWarp.CMYKAt(x, y)
			cd := dWarp.CMYKAt(x, y)
			c := color.CMYK{
				C: avgU8(cs.C, cd.C, t),
				M: avgU8(cs.M, cd.M, t),
				Y: avgU8(cs.Y, cd.Y, t),
				K: avgU8(cs.K, cd.K, t),
			}
			img.SetCMYK(x, y, c)
		}
	}
	return img, nil
}


// Morph morphs one image to another by warping an input mesh some fraction of
// the way to an output mesh.
func Morph(sImg, dImg image.Image, sMesh, dMesh *Mesh, t float64) (image.Image, error) {
	if reflect.TypeOf(sImg) != reflect.TypeOf(dImg) {
		panic(fmt.Sprintf("morphing from %T to %T is not yet implemented", sImg, dImg)) // TODO: implement
	}
	switch sImg.(type) {
	case *image.NRGBA:
		return morphNRGBA(sImg.(*image.NRGBA), dImg.(*image.NRGBA),
			sMesh, dMesh, t)
	case *image.CMYK:
		return morphCMYK(sImg.(*image.CMYK), dImg.(*image.CMYK),
			sMesh, dMesh, t)
	default:
		panic(fmt.Sprintf("morphing from %T to %T is not yet implemented", sImg, dImg)) // TODO: implement
	}
}
