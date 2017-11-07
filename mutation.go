package main

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/aurelien-rainone/evolve/operators"
)

type mutationOptions struct {
	addPolygonMutation    number.ProbabilityGenerator
	removePolygonMutation number.ProbabilityGenerator
	swapPolygonsMutation  number.ProbabilityGenerator
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

	if op.options.addPolygonMutation.NextValue().NextEvent(rng) {
		// add a new random polygon
		img.polys = append(img.polys,
			randomPoly(img, newPolyMinPoints, newPolyMaxPoints, rng))
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

	// returns cloned image, possibily mutated
	return img
}
