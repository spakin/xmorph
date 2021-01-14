// This file provides functions for distorting images.

package morph

/*
#include <xmorph/mesh.h>
#include <xmorph/mesh_t.h>
#include <xmorph/warp2.h>
*/
import "C"
import (
	"fmt"
	"image"
)

// warpNRGBA warps an NRGBA image.
func warpNRGBA(img *image.NRGBA, src, dst *Mesh) *image.NRGBA {
	bnds := img.Bounds()
	wd := bnds.Max.X - bnds.Min.X
	ht := bnds.Max.Y - bnds.Min.Y
	out := make([]uint8, len(img.Pix))
	C.warp_image_versatile(
		// Source information
		(*C.PIXEL_TYPE)(&img.Pix[0]),
		C.int(wd), C.int(ht), 4,
		C.int(img.Stride), 4,
		// Destination information
		(*C.PIXEL_TYPE)(&out[0]),
		C.int(wd), C.int(ht), 4,
		C.int(img.Stride), 4,
		// Mesh information
		src.mesh.x, src.mesh.y,
		dst.mesh.x, dst.mesh.y,
		C.int(src.mesh.nx), C.int(src.mesh.ny))
	return &image.NRGBA{
		Pix:    out,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

// Warp distorts an image by warping an input mesh to an output mesh.
func Warp(img image.Image, src, dst *Mesh) (image.Image, error) {
	if C.meshCompatibilityCheck(src.mesh, dst.mesh) != 0 {
		return nil, fmt.Errorf("incompatible meshes passed to InterpolateMeshes")
	}
	switch img := img.(type) {
	case *image.NRGBA:
		return warpNRGBA(img, src, dst), nil
	default:
		return nil, fmt.Errorf("warping of %T image types is not yet implemented", img)
	}	
}
