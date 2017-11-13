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

	// set image-level mutations
	if prob, err = number.NewProbability(appConfig.Mutation.Image.AddPoly); err != nil {
		return nil, fmt.Errorf("add-polygon mutation rate error: %v", err)
	}
	mutater.addPolygonMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Image.RemovePoly); err != nil {
		return nil, fmt.Errorf("remove-polygon mutation rate error: %v", err)
	}
	mutater.removePolygonMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Image.SwapPolys); err != nil {
		return nil, fmt.Errorf("swap-polygon mutation rate error: %v", err)
	}
	mutater.swapPolygonsMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Image.BackgroundColor); err != nil {
		return nil, fmt.Errorf("background-color mutation rate error: %v", err)
	}
	mutater.backgroundColorMutation = number.NewConstantProbabilityGenerator(prob)

	// set polygon-level mutations
	if prob, err = number.NewProbability(appConfig.Mutation.Polygon.AddPoint); err != nil {
		return nil, fmt.Errorf("add-point mutation rate error: %v", err)
	}
	mutater.addPointMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Polygon.RemovePoint); err != nil {
		return nil, fmt.Errorf("remove-point mutation rate error: %v", err)
	}
	mutater.removePointMutation = number.NewConstantProbabilityGenerator(prob)

	if prob, err = number.NewProbability(appConfig.Mutation.Polygon.ChangeColor); err != nil {
		return nil, fmt.Errorf("change-polygon-color mutation rate error: %v", err)
	}
	mutater.changePolyColorMutation = number.NewConstantProbabilityGenerator(prob)

	// set point-level mutations
	if prob, err = number.NewProbability(appConfig.Mutation.Point.Move); err != nil {
		return nil, fmt.Errorf("move-point mutation rate error: %v", err)
	}
	mutater.movePointMutation = number.NewConstantProbabilityGenerator(prob)

	impl, err := operators.NewAbstractMutation(mutater)
	mutater.impl = impl
	return impl, err
}

type imageDNAMutater struct {
	impl *operators.AbstractMutation

	// image-level mutations
	addPolygonMutation      number.ProbabilityGenerator
	removePolygonMutation   number.ProbabilityGenerator
	swapPolygonsMutation    number.ProbabilityGenerator
	backgroundColorMutation number.ProbabilityGenerator

	// polygon-level mutations
	addPointMutation        number.ProbabilityGenerator
	removePointMutation     number.ProbabilityGenerator
	changePolyColorMutation number.ProbabilityGenerator

	// point-level mutations
	movePointMutation number.ProbabilityGenerator
}

func (op *imageDNAMutater) Mutate(c framework.Candidate, rng *rand.Rand) framework.Candidate {
	// mutates a copy of the image, mutation do not touch the original
	img := c.(*imageDNA).clone()

	if op.addPolygonMutation.NextValue().NextEvent(rng) {
		if len(img.polys) < appConfig.Image.MaxPolys {
			// add a new random polygon
			img.polys = append(img.polys,
				randomPoly(img, appConfig.Polygon.MinPoints, appConfig.Polygon.MaxPoints, rng))
		}
	}

	if op.removePolygonMutation.NextValue().NextEvent(rng) {
		if len(img.polys) > appConfig.Image.MinPolys {
			// find removal index
			idx := rng.Intn(len(img.polys))
			// split slice before and after, and append those 2 parts together
			img.polys = append(img.polys[:idx], img.polys[idx+1:]...)
		}
	}

	if op.swapPolygonsMutation.NextValue().NextEvent(rng) {
		// swap 2 random polygons
		idx1, idx2 := rng.Intn(len(img.polys)), rng.Intn(len(img.polys))
		img.polys[idx1], img.polys[idx2] = img.polys[idx2], img.polys[idx1]
	}

	if op.backgroundColorMutation.NextValue().NextEvent(rng) {
		img.bck = randomColorNoAlpha(rng)
	}

	for i := 0; i < len(img.polys); i++ {
		poly := &img.polys[i]

		if op.changePolyColorMutation.NextValue().NextEvent(rng) {
			// change poly color
			// TODO: which is best? try to evolve current color or start with a
			// random one
			poly.col = randomColor(rng)
			//evolveColor(&poly.col, rng)
		}

		if op.addPointMutation.NextValue().NextEvent(rng) {
			numPts := len(poly.pts)
			if numPts < appConfig.Polygon.MaxPoints {
				// find insertion index
				idx := 1 + rng.Intn(numPts-1)
				// insert point at the middle of prev and next points
				poly.insert(idx, poly.pts[idx-1].Add(poly.pts[idx]).Div(2))
			}
		}

		if op.removePointMutation.NextValue().NextEvent(rng) {
			numPts := len(poly.pts)
			if numPts > appConfig.Polygon.MinPoints {
				// find removal index
				idx := rng.Intn(numPts)
				// split slice before and after, and append those 2 parts together
				poly.pts = append(poly.pts[:idx], poly.pts[idx+1:]...)
			}
		}

		if op.removePointMutation.NextValue().NextEvent(rng) {
			numPts := len(poly.pts)
			if numPts > appConfig.Polygon.MinPoints {
				// find removal index
				idx := rng.Intn(numPts)
				// split slice before and after, and append those 2 parts together
				poly.pts = append(poly.pts[:idx], poly.pts[idx+1:]...)
			}
		}

		for j := 0; j < len(poly.pts); j++ {
			//pt := &poly.pts[j]
			if op.movePointMutation.NextValue().NextEvent(rng) {
				// TODO: compute margin
				poly.pts[j] = randomPoint(img, 10, rng)
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
