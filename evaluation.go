package main

import (
	"image"

	"github.com/aurelien-rainone/evolve/framework"
)

type fitnessEvaluator struct {
	img image.Image // reference image
}

func (fe *fitnessEvaluator) Fitness(c framework.Candidate, pop []framework.Candidate) float64 {
	//img := c.(*imageDNA)
	return 0
}

func (fe *fitnessEvaluator) IsNatural() bool {
	return true
}
