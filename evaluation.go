package main

import (
	"image"

	"github.com/aurelien-rainone/evolve/framework"
)

type fitnessEvaluator struct {
	img *image.RGBA // reference image
}

// compare a reference image to a test image and returns the difference
func imageDiff(ref, img *image.RGBA) float64 {
	b := ref.Bounds()
	w, h := b.Dx(), b.Dy()
	// create a data set for storing all diff errors
	data := make([]float64, w*h)
	var offset int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			offset = y*ref.Stride + x*4
			r := ref.Pix[offset+0] - img.Pix[offset+0]
			g := ref.Pix[offset+1] - img.Pix[offset+1]
			b := ref.Pix[offset+2] - img.Pix[offset+2]
			a := ref.Pix[offset+3] - img.Pix[offset+3]
			data[y*w+x] = float64(r + g + b + a)
		}
	}
	ds := framework.NewDataSet(framework.WithPrePopulatedDataSet(data))
	return ds.Variance()
}

func (fe *fitnessEvaluator) Fitness(c framework.Candidate, pop []framework.Candidate) float64 {
	return imageDiff(fe.img, c.(*imageDNA).render())
}

func (fe *fitnessEvaluator) IsNatural() bool {
	// the lesser the fitness the better
	return false
}
