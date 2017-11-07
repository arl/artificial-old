package main

import (
	"image"
	"image/color"
)

// poly represents a polygon of the image
type poly struct {
	col color.RGBA
	pts []image.Point
}

// imageDNA is a gene coding for an image made of polygons
type imageDNA struct {
	w, h  int
	polys []poly
}

// clone returns a new imageDNA that is an exact copy of the receiver
func (img *imageDNA) clone() *imageDNA {
	// copy polygon slice
	polys := make([]poly, 0, len(img.polys))
	for _, p := range img.polys {
		poly := poly{col: p.col}
		// copy points slice
		poly.pts = make([]image.Point, 0, len(p.pts))
		copy(poly.pts, p.pts)
		polys = append(polys, poly)
	}
	return &imageDNA{polys: polys, w: img.w, h: img.h}
}
