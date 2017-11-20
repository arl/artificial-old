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
		img            = c.(*imageDNA).render() // rendered chromosome
		b              = fe.img.Bounds()        // image bounds
		w, h           = b.Dx(), b.Dy()
		data           = make([]float64, w*h) // data to store all pixel differences
		off            int
		diff           int64
		rr, rg, rb, ra uint8
		ir, ig, ib, ia uint8
	)

	// compare a reference image to a test image and returns the difference
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			off = y*fe.img.Stride + x*4
			rr, rg, rb, ra = fe.img.Pix[off+0], fe.img.Pix[off+1], fe.img.Pix[off+2], fe.img.Pix[off+3]
			ir, ig, ib, ia = img.Pix[off+0], img.Pix[off+1], img.Pix[off+2], img.Pix[off+3]
			diff = abs(int64(rr) - int64(ir))
			diff += abs(int64(rg) - int64(ig))
			diff += abs(int64(rb) - int64(ib))
			diff += abs(int64(ra) - int64(ia))
			data[y*w+x] = float64(diff)
		}
	}

	// TODO: if we are only going to do the arithmetic mean, no need to save
	// each value, just sum them and divide by the number, avoiding costing
	// allocation
	ds := framework.NewDataSet(framework.WithPrePopulatedDataSet(data))
	return ds.ArithmeticMean()
}

func (fe *fitnessEvaluator) IsNatural() bool {
	// the lesser the fitness the better
	return false
}
