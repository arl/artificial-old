package main

import (
	"image"
	"image/color"
)

// Poly represents a polygon of the image
type Poly struct {
	pts []image.Point
	col color.RGBA
}

// ImageDNA is a gene coding for an image made of polygons
type ImageDNA struct {
	polys []Poly
	w, h  int
}
