// This file provides functions for morphing images.

package xmorph

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

// morphNRGBA morphs two NRGBA images.
func morphNRGBA(sImg, dImg *image.NRGBA, sMesh, dMesh *Mesh, t float64) (*image.NRGBA, error) {
	// Create an mesh intermediate to the source and destination meshes.
	mMesh, err := InterpolateMeshes(sMesh, dMesh, t)
	if err != nil {
		return nil, err
	}

	// Separately warp the source and destination images to the
	// intermediate mesh.
	sw, err := Warp(sImg, sMesh, mMesh, 1.0)
	if err != nil {
		return nil, err
	}
	dw, err := Warp(dImg, dMesh, mMesh, 1.0)
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

// morphGray morphs two Gray images.
func morphGray(sImg, dImg *image.Gray, sMesh, dMesh *Mesh, t float64) (*image.Gray, error) {
	// Separately warp the source and destination images.
	sw, err := Warp(sImg, sMesh, dMesh, t)
	if err != nil {
		return nil, err
	}
	dw, err := Warp(dImg, sMesh, dMesh, 1.0-t)
	if err != nil {
		return nil, err
	}
	sWarp, dWarp := sw.(*image.Gray), dw.(*image.Gray)

	// Perform a weighted average of the source and destination images'
	// colors to produce a final image.
	bnds := sWarp.Bounds()
	img := image.NewGray(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			cs := sWarp.GrayAt(x, y)
			cd := dWarp.GrayAt(x, y)
			c := color.Gray{
				Y: avgU8(cs.Y, cd.Y, t),
			}
			img.SetGray(x, y, c)
		}
	}
	return img, nil
}

// morphCMYK morphs two CMYK images.
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

// morphAlpha morphs two Alpha images.
func morphAlpha(sImg, dImg *image.Alpha, sMesh, dMesh *Mesh, t float64) (*image.Alpha, error) {
	// Separately warp the source and destination images.
	sw, err := Warp(sImg, sMesh, dMesh, t)
	if err != nil {
		return nil, err
	}
	dw, err := Warp(dImg, sMesh, dMesh, 1.0-t)
	if err != nil {
		return nil, err
	}
	sWarp, dWarp := sw.(*image.Alpha), dw.(*image.Alpha)

	// Perform a weighted average of the source and destination images'
	// colors to produce a final image.
	bnds := sWarp.Bounds()
	img := image.NewAlpha(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			cs := sWarp.AlphaAt(x, y)
			cd := dWarp.AlphaAt(x, y)
			c := color.Alpha{
				A: avgU8(cs.A, cd.A, t),
			}
			img.SetAlpha(x, y, c)
		}
	}
	return img, nil
}

// morphAny morphs two images of any type by first converting them to NRGBA and
// then invoking morphNRGBA.
func morphAny(sImg, dImg image.Image, sMesh, dMesh *Mesh, t float64) (*image.NRGBA, error) {
	bnds := sImg.Bounds()
	sNrgba := image.NewNRGBA(bnds)
	dNrgba := image.NewNRGBA(bnds)
	cm := sNrgba.ColorModel()
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			sc := sImg.At(x, y)
			sNrgba.Set(x, y, cm.Convert(sc))
			dc := dImg.At(x, y)
			dNrgba.Set(x, y, cm.Convert(dc))
		}
	}
	return morphNRGBA(sNrgba, dNrgba, sMesh, dMesh, t)
}

// Morph morphs one image to another by warping a source mesh some fraction of
// the way to a destination mesh.
func Morph(sImg, dImg image.Image, sMesh, dMesh *Mesh, t float64) (image.Image, error) {
	if sImg.Bounds() != dImg.Bounds() {
		return nil, fmt.Errorf("images to morph must have the same bounds")
	}
	if reflect.TypeOf(sImg) != reflect.TypeOf(dImg) {
		return morphAny(sImg, dImg, sMesh, dMesh, t)
	}
	switch sImg.(type) {
	case *image.NRGBA:
		return morphNRGBA(sImg.(*image.NRGBA), dImg.(*image.NRGBA),
			sMesh, dMesh, t)
	case *image.Gray:
		return morphGray(sImg.(*image.Gray), dImg.(*image.Gray),
			sMesh, dMesh, t)
	case *image.CMYK:
		return morphCMYK(sImg.(*image.CMYK), dImg.(*image.CMYK),
			sMesh, dMesh, t)
	case *image.Alpha:
		return morphAlpha(sImg.(*image.Alpha), dImg.(*image.Alpha),
			sMesh, dMesh, t)
	default:
		return morphAny(sImg, dImg, sMesh, dMesh, t)
	}
}
