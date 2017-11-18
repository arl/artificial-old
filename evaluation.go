package main

import (
	"image"
	"image/color"

	"github.com/aurelien-rainone/evolve/framework"
)

type fitnessEvaluator struct {
	img *image.RGBA // reference image
}

// compare a reference image to a test image and returns the difference

// TODO: look at https://github.com/mapbox/pixelmatch/blob/master/index.js to
// see preceived differences in color

func imageDiff(ref, img *image.RGBA) float64 {
	b := ref.Bounds()
	w, h := b.Dx(), b.Dy()
	// create a data set for storing all diff errors
	data := make([]float64, w*h)
	//var offset int
	var diff int64
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			//r1, g1, b1, _, r2, g2, b2, _ := ref.At(x, y).RGBA(), img.At(x, y).RGBA()
			diff = diffColor(ref.At(x, y), img.At(x, y))
			//offset = y*ref.Stride + x*4
			//r := ref.Pix[offset+0] - img.Pix[offset+0]
			//g := ref.Pix[offset+1] - img.Pix[offset+1]
			//b := ref.Pix[offset+2] - img.Pix[offset+2]
			// alpha difference do not count
			//a := ref.Pix[offset+3] - img.Pix[offset+3]
			data[y*w+x] = float64(diff)
		}
	}
	ds := framework.NewDataSet(framework.WithPrePopulatedDataSet(data))
	return ds.ArithmeticMean()
}

func diffColor(c1, c2 color.Color) int64 {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	var diff int64
	diff += abs(int64(r1) - int64(r2))
	diff += abs(int64(g1) - int64(g2))
	diff += abs(int64(b1) - int64(b2))
	diff += abs(int64(a1) - int64(a2))
	return diff
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func (fe *fitnessEvaluator) Fitness(c framework.Candidate, pop []framework.Candidate) float64 {
	return imageDiff(fe.img, c.(*imageDNA).render())
}

func (fe *fitnessEvaluator) IsNatural() bool {
	// the lesser the fitness the better
	return false
}
