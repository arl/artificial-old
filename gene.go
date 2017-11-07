package main

import (
	"image"
	"image/color"
)

// poly represents a polygon of the image
type poly struct {
	pts []image.Point
	col color.RGBA
}

// imageDNA is a gene coding for an image made of polygons
type imageDNA struct {
	polys []poly
	w, h  int
}
