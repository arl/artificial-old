package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/llgcode/draw2d/draw2dimg"
)

// poly represents a polygon of the image
type poly struct {
	col color.RGBA
	pts []image.Point
}

func (p *poly) insert(idx int, pt image.Point) {
	// append a zero-value at the back
	p.pts = append(p.pts, image.Point{})
	// right-shift all elements after the insertion point
	copy(p.pts[idx+1:], p.pts[idx:])
	// set the inserted element at given index
	p.pts[idx] = pt
}

// imageDNA is a gene coding for an image made of polygons
type imageDNA struct {
	w, h  int
	polys []poly
}

// clone returns a new imageDNA that is an exact copy of the receiver
func (img *imageDNA) clone() *imageDNA {
	// copy polygon slice
	polys := make([]poly, len(img.polys))
	for i, p := range img.polys {
		poly := poly{col: p.col}
		// copy points slice
		poly.pts = make([]image.Point, len(p.pts))
		copy(poly.pts, p.pts)
		polys[i] = poly
	}
	return &imageDNA{polys: polys, w: img.w, h: img.h}
}

func (img *imageDNA) render() *image.RGBA {
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, img.w, img.h))
	gc := draw2dimg.NewGraphicContext(dest)
	gc.SetLineWidth(1)

	for i := 0; i < len(img.polys); i++ {
		poly := img.polys[i]
		// set brush color
		gc.SetFillColor(poly.col)
		gc.SetStrokeColor(poly.col)

		// draw polygon as a closed path
		gc.MoveTo(float64(poly.pts[0].X), float64(poly.pts[0].Y))
		for j := 1; j < len(poly.pts); j++ {
			pt := poly.pts[j]
			gc.LineTo(float64(pt.X), float64(pt.Y))
		}
		gc.SetLineWidth(0)
		gc.Close()
		gc.FillStroke()
	}
	return dest
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
