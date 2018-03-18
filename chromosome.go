package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	gorng "github.com/leesper/go_rng"
)

// poly represents a polygon of the image
type poly struct {
	col color.Color
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

// randomSimplePoly creates and returns a random simple polygon.
func randomSimplePoly(img *imageDNA, minPts, maxPts int, rng *rand.Rand) poly {
	poly := poly{}

	// create random number of points
	numPts := minPts + rng.Intn(maxPts-minPts)

	// compute random polygon average radius (5-30% of the image size)
	minRadius := (img.w * 5) / 100
	maxRadius := (img.w * 30) / 100
	margin := minRadius + rng.Intn(maxRadius-minRadius)

	// random point to be the polygon center
	center := randomPoint(img, margin, rng)

	// use polygon generator
	poly.pts = generatePolygon(center, float64(margin), 0.7, 0.5, numPts, rng)

	// set random color
	poly.col = randomColor(rng)
	return poly
}

// randomPoly creates and returns a random polygon.
func randomPoly(img *imageDNA, minPts, maxPts int, rng *rand.Rand) poly {
	poly := poly{}
	// create random number of points
	var numPts int
	if maxPts == minPts {
		numPts = maxPts
	} else {
		numPts = minPts + rng.Intn(maxPts-minPts)
	}
	poly.pts = make([]image.Point, numPts)
	for j := 0; j < numPts; j++ {
		// each point is random
		poly.pts[j] = randomPoint(img, 0, rng)
	}
	// set random color
	poly.col = randomColor(rng)
	return poly
}

// randomPoint creates and returns a random point in the image
//
// margin is the min distance in pixel from the image border
func randomPoint(img *imageDNA, margin int, rng *rand.Rand) image.Point {
	return image.Point{
		X: margin + rng.Intn(img.w-2*margin),
		Y: margin + rng.Intn(img.h-2*margin),
	}
}

// randomPoint returns a random color
func randomColor(rng *rand.Rand) color.Color {
	return color.NRGBA{
		R: byte(rng.Intn(255)),
		G: byte(rng.Intn(255)),
		B: byte(rng.Intn(255)),
		A: byte(10 + rng.Intn(50)),
	}
}

var gauss *gorng.GaussianGenerator

func init() {
	// instantiate the gaussian generator
	gauss = gorng.NewGaussianGenerator(time.Now().UnixNano())
}

// Start with the centre of the polygon at ctrX, ctrY,
// then creates the polygon by sampling points on a circle around the centre.
// Randon noise is added by varying the angular spacing between sequential points,
// and by varying the radial distance of each point from the centre.
//
// Params:
// ctrX, ctrY - coordinates of the "centre" of the polygon
// aveRadius - in px, the average radius of this polygon, this roughly controls how large the polygon is, really only useful for order of magnitude.
// irregularity - [0,1] indicating how much variance there is in the angular spacing of vertices. [0,1] will map to [0, 2pi/numberOfVerts]
// spikeyness - [0,1] indicating how much variance there is in each vertex from the circle of radius aveRadius. [0,1] will map to [0, aveRadius]
// numPts - self-explanatory
//
// Returns a list of vertices, in CCW order.
// Taken from:
// https://stackoverflow.com/questions/8997099/algorithm-to-generate-random-2d-polygon
func generatePolygon(ctr image.Point, avgRadius, irregularity, spikeyness float64, numPts int, rng *rand.Rand) []image.Point {
	irregularity = f64Clip(irregularity, 0, 1) * 2 * math.Pi / float64(numPts)
	spikeyness = f64Clip(spikeyness, 0, 1) * avgRadius

	// generate n angle steps
	angleSteps := make([]float64, numPts)
	lower := (2 * math.Pi / float64(numPts)) - irregularity
	upper := (2 * math.Pi / float64(numPts)) + irregularity
	var sum float64
	for i := 0; i < numPts; i++ {
		tmp := lower + rng.Float64()*(upper-lower)
		angleSteps[i] = tmp
		sum = sum + tmp
	}

	// normalize the steps so that point 0 and point n+1 are the same
	k := sum / (2 * math.Pi)
	for i := 0; i < numPts; i++ {
		angleSteps[i] = angleSteps[i] / k
	}

	// now generate the points
	points := make([]image.Point, numPts)
	angle := rng.Float64() * 2 * math.Pi
	var x, y float64
	for i := 0; i < numPts; i++ {
		ri := f64Clip(gauss.Gaussian(avgRadius, spikeyness), 0, 2*avgRadius)
		x = float64(ctr.X) + ri*math.Cos(angle)
		y = float64(ctr.Y) + ri*math.Sin(angle)
		points[i] = image.Pt(int(x), int(y))
		angle = angle + angleSteps[i]
	}
	return points
}

func f64Clip(x, min, max float64) float64 {
	if min > max {
		return x
	}
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
