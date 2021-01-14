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

// warpUint8Slice warps any image type that's representable as a slice of
// alternating channel values, each of which is of type uint8.
func warpUint8Slice(pix []uint8, ystr, nchan int, bnds image.Rectangle, src, dst *Mesh) []uint8 {
	wd := bnds.Max.X - bnds.Min.X
	ht := bnds.Max.Y - bnds.Min.Y
	out := make([]uint8, len(pix))
	C.warp_image_versatile(
		// Source information
		(*C.PIXEL_TYPE)(&pix[0]),
		C.int(wd), C.int(ht), C.int(nchan), C.int(ystr), C.int(nchan),
		// Destination information
		(*C.PIXEL_TYPE)(&out[0]),
		C.int(wd), C.int(ht), C.int(nchan), C.int(ystr), C.int(nchan),
		// Mesh information
		src.mesh.x, src.mesh.y,
		dst.mesh.x, dst.mesh.y,
		C.int(src.mesh.nx), C.int(src.mesh.ny))
	return out
}

// warpAlpha warps an Alpha image.
func warpAlpha(img *image.Alpha, src, dst *Mesh) *image.Alpha {
	out := warpUint8Slice(img.Pix, img.Stride, 1, img.Rect, src, dst)
	return &image.Alpha{
		Pix:    out,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

// warpNRGBA warps an NRGBA image.
func warpNRGBA(img *image.NRGBA, src, dst *Mesh) *image.NRGBA {
	out := warpUint8Slice(img.Pix, img.Stride, 4, img.Rect, src, dst)
	return &image.NRGBA{
		Pix:    out,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

// warpCMYK warps a CMYK image.
func warpCMYK(img *image.CMYK, src, dst *Mesh) *image.CMYK {
	out := warpUint8Slice(img.Pix, img.Stride, 4, img.Rect, src, dst)
	return &image.CMYK{
		Pix:    out,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

// warpGray warps a Gray image.
func warpGray(img *image.Gray, src, dst *Mesh) *image.Gray {
	out := warpUint8Slice(img.Pix, img.Stride, 1, img.Rect, src, dst)
	return &image.Gray{
		Pix:    out,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
}

// warpRGBA warps an RGBA image.
func warpRGBA(img *image.RGBA, src, dst *Mesh) *image.RGBA {
	out := warpUint8Slice(img.Pix, img.Stride, 4, img.Rect, src, dst)
	return &image.RGBA{
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
	case *image.RGBA:
		return warpRGBA(img, src, dst), nil
	case *image.NRGBA:
		return warpNRGBA(img, src, dst), nil
	case *image.Gray:
		return warpGray(img, src, dst), nil
	case *image.CMYK:
		return warpCMYK(img, src, dst), nil
	case *image.Alpha:
		return warpAlpha(img, src, dst), nil
	default:
		return nil, fmt.Errorf("warping of %T image types is not yet implemented", img)
	}
}
