package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/aurelien-rainone/evolve/operators"
)

type mutationOptions struct {
	addPolygonMutation      number.ProbabilityGenerator
	removePolygonMutation   number.ProbabilityGenerator
	swapPolygonsMutation    number.ProbabilityGenerator
	changePolyColorMutation number.ProbabilityGenerator
}

func newImageDNAMutation(options mutationOptions) (*operators.AbstractMutation, error) {
	// set default mutation values
	if options.addPolygonMutation == nil {
		options.addPolygonMutation = number.NewConstantProbabilityGenerator(number.ProbabilityZero)
	}
	if options.removePolygonMutation == nil {
		options.removePolygonMutation = number.NewConstantProbabilityGenerator(number.ProbabilityZero)
	}
	if options.swapPolygonsMutation == nil {
		options.swapPolygonsMutation = number.NewConstantProbabilityGenerator(number.ProbabilityZero)
	}
	if options.changePolyColorMutation == nil {
		options.changePolyColorMutation = number.NewConstantProbabilityGenerator(number.ProbabilityZero)
	}

	mutater := &imageDNAMutater{options: options}
	impl, err := operators.NewAbstractMutation(mutater)
	mutater.impl = impl
	return impl, err
}

type imageDNAMutater struct {
	impl    *operators.AbstractMutation
	options mutationOptions
}

func (op *imageDNAMutater) Mutate(c framework.Candidate, rng *rand.Rand) framework.Candidate {
	// mutates a copy of the image, mutation do not touch the original
	img := c.(*imageDNA).clone()

	// image-level mutations

	if op.options.addPolygonMutation.NextValue().NextEvent(rng) {
		// add a new random polygon
		img.polys = append(img.polys,
			randomPoly(img, appConfig.Polygon.MinPoints, appConfig.Polygon.MaxPoints, rng))
	}

	if op.options.removePolygonMutation.NextValue().NextEvent(rng) {
		// remove random polygon
		idx := rng.Intn(len(img.polys))
		img.polys = append(img.polys[:idx], img.polys[idx+1:]...)
	}

	if op.options.swapPolygonsMutation.NextValue().NextEvent(rng) {
		// swap 2 random polygons
		idx1, idx2 := rng.Intn(len(img.polys)), rng.Intn(len(img.polys))
		img.polys[idx1], img.polys[idx2] = img.polys[idx2], img.polys[idx1]
	}

	// polygon-level mutations

	for i := 0; i < len(img.polys); i++ {
		poly := &img.polys[i]

		if op.options.changePolyColorMutation.NextValue().NextEvent(rng) {
			// change poly color
			evolveColor(&poly.col, rng)
		}
	}

	// returns cloned image, possibily mutated
	return img
}

const (
	// max percentage of decrease/increase in value of a color component
	maxByteEvolutionPercent = 10
	maxByteEvolution        = math.MaxUint8 / maxByteEvolutionPercent
)

func evolveColor(c *color.RGBA, rng *rand.Rand) {
	evolveByte := func(b byte) byte {
		// max byte value
		var maxVal byte = math.MaxUint8
		if b < math.MaxUint8-maxByteEvolution {
			maxVal = b + maxByteEvolution
		}

		// min byte value
		var minVal byte = 0
		if b > maxByteEvolution {
			minVal = b - maxByteEvolution
		}
		return minVal + byte(rng.Intn(int(maxVal-minVal)))
	}

	// we want each color component to get +/- 10% than their current value
	c.R = evolveByte(c.R)
	c.G = evolveByte(c.G)
	c.B = evolveByte(c.B)
	c.A = evolveByte(c.A)
}
