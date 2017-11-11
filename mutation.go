package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/aurelien-rainone/evolve/operators"
)

func newImageDNAMutation() (*operators.AbstractMutation, error) {
	// create and configure mutater with all mutation rates
	mutater := &imageDNAMutater{}

	var (
		prob number.Probability
		err  error
	)

	// set polygon mutations
	if prob, err = number.NewProbability(appConfig.Mutation.Polygon.Add); err != nil {
		return nil, fmt.Errorf("add-polygon mutation rate error: %v", err)
	}
	mutater.addPolygonMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Polygon.Remove); err != nil {
		return nil, fmt.Errorf("remove-polygon mutation rate error: %v", err)
	}
	mutater.removePolygonMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Polygon.Swap); err != nil {
		return nil, fmt.Errorf("swap-polygon mutation rate error: %v", err)
	}
	mutater.swapPolygonsMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Polygon.Color); err != nil {
		return nil, fmt.Errorf("change-polygon-color mutation rate error: %v", err)
	}
	mutater.changePolyColorMutation = number.NewConstantProbabilityGenerator(prob)

	// set point mutations
	if prob, err = number.NewProbability(appConfig.Mutation.Point.Add); err != nil {
		return nil, fmt.Errorf("add-point mutation rate error: %v", err)
	}
	mutater.addPointMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Point.Remove); err != nil {
		return nil, fmt.Errorf("remove-point mutation rate error: %v", err)
	}
	mutater.removePointMutation = number.NewConstantProbabilityGenerator(prob)

	impl, err := operators.NewAbstractMutation(mutater)
	mutater.impl = impl
	return impl, err
}

type imageDNAMutater struct {
	impl *operators.AbstractMutation

	// polygon mutations
	addPolygonMutation      number.ProbabilityGenerator
	removePolygonMutation   number.ProbabilityGenerator
	swapPolygonsMutation    number.ProbabilityGenerator
	changePolyColorMutation number.ProbabilityGenerator

	// point mutations
	addPointMutation    number.ProbabilityGenerator
	removePointMutation number.ProbabilityGenerator
}

func (op *imageDNAMutater) Mutate(c framework.Candidate, rng *rand.Rand) framework.Candidate {
	// mutates a copy of the image, mutation do not touch the original
	img := c.(*imageDNA).clone()

	if op.addPolygonMutation.NextValue().NextEvent(rng) {
		// add a new random polygon
		img.polys = append(img.polys,
			randomPoly(img, appConfig.Polygon.MinPoints, appConfig.Polygon.MaxPoints, rng))
	}

	if op.removePolygonMutation.NextValue().NextEvent(rng) {
		// remove random polygon
		idx := rng.Intn(len(img.polys))
		img.polys = append(img.polys[:idx], img.polys[idx+1:]...)
	}

	if op.swapPolygonsMutation.NextValue().NextEvent(rng) {
		// swap 2 random polygons
		idx1, idx2 := rng.Intn(len(img.polys)), rng.Intn(len(img.polys))
		img.polys[idx1], img.polys[idx2] = img.polys[idx2], img.polys[idx1]
	}

	for i := 0; i < len(img.polys); i++ {
		poly := &img.polys[i]
		numPts := len(poly.pts)

		if op.changePolyColorMutation.NextValue().NextEvent(rng) {
			// change poly color
			// TODO: which is best? try to evolve current color or start with a
			// random one
			poly.col = randomColor(rng)
			//evolveColor(&poly.col, rng)
		}

		if op.addPointMutation.NextValue().NextEvent(rng) {
			if numPts < appConfig.Polygon.MaxPoints {
				// find insertion index
				idx := 1 + rng.Intn(numPts-1)
				// insert point at the middle of prev and next points
				poly.insert(idx, poly.pts[idx-1].Add(poly.pts[idx]).Div(2))
			}
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
