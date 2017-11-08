package main

import (
	"image"
	"image/color"
	"math/rand"
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

// randomPoint creates and returns a random polygon made of points in the image,
// with minPts < numPts < maxPts
func randomPoly(img *imageDNA, minPts, maxPts int, rng *rand.Rand) poly {
	poly := poly{}
	// create random number of points
	numPoints := minPts + rng.Intn(maxPts-minPts)
	poly.pts = make([]image.Point, numPoints)
	for j := 0; j < numPoints; j++ {
		// each point is random
		poly.pts[j] = randomPoint(img, rng)
	}
	// set random color
	poly.col = randomColor(rng)
	return poly
}

// randomPoint creates and returns a random point in the image
func randomPoint(img *imageDNA, rng *rand.Rand) image.Point {
	return image.Pt(rng.Intn(img.w), rng.Intn(img.w))
}

// randomPoint returns a random RGBA color
func randomColor(rng *rand.Rand) color.RGBA {
	return color.RGBA{
		R: byte(rng.Intn(255)),
		G: byte(rng.Intn(255)),
		B: byte(rng.Intn(255)),
		A: byte(10 + rng.Intn(50)),
	}
}
