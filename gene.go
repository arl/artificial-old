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
