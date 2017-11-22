package main

import (
	"image"

	"github.com/aurelien-rainone/evolve/framework"
)

type fitnessEvaluator struct {
	img *image.RGBA // reference image
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func (fe *fitnessEvaluator) Fitness(c framework.Candidate, pop []framework.Candidate) float64 {
	var (
		img  = c.(*imageDNA).render() // rendered chromosome
		b    = fe.img.Bounds()        // image bounds
		w, h = b.Dx(), b.Dy()
		off  int
		diff int64
	)

	// compare a reference image to a test image and returns the difference
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			off = y*fe.img.Stride + x*4
			diff += abs(int64(fe.img.Pix[off+0])+int64(fe.img.Pix[off+1])+int64(fe.img.Pix[off+2])) -
				abs(int64(img.Pix[off+0])+int64(img.Pix[off+1])+int64(img.Pix[off+2]))
		}
	}
	return float64(diff)
}

func (fe *fitnessEvaluator) IsNatural() bool {
	// the lesser the fitness the better
	return false
}
