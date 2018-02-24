package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/fogleman/gg"
)

func testRandomPoly(img *imageDNA, numPts int, rng *rand.Rand) poly {
	poly := poly{}

	// compute random polygon average radius (5-30% of the image size)
	minRadius := (img.w * 5) / 100
	maxRadius := (img.w * 30) / 100
	margin := minRadius + rng.Intn(maxRadius-minRadius)
	fmt.Printf("image: %vx%v, margin:%v\n", img.w, img.h, margin)
	// random center
	center := randomPoint(img, margin, rng)

	// use polygon generator
	poly.pts = generatePolygon(center, float64(margin), 0.7, 0.5, numPts, rng)
	// set random color
	poly.col = randomColor(rng)
	return poly
}

func TestPolygonGeneration(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dna := imageDNA{w: 500, h: 500}
	dna.polys = []poly{testRandomPoly(&dna, 6, rng)}
	img := dna.render()
	dc := gg.NewContextForRGBA(img)
	dc.SavePNG("poly.png")
}
