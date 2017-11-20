package main

import (
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/operators"
)

func newImageDNACrossover(options ...operators.Option) (*operators.AbstractCrossover, error) {
	// TODO: for now hardcoded options
	return operators.NewAbstractCrossover(imageDNAMater{}, operators.ConstantCrossoverPoints(1))
}

type imageDNAMater struct{}

// imageDNAMater implements an unequal crossover as chromosomes (imageDNA
// instances) may code for images with different number of polygons.
// In order to reduce the problem to an equal-length crossover, we only consider
// a sub-sequence of the longest chromosome, that has the same length than the
// shorter one. The start index of the considered subsequence is randomly chosen
// between 0 and the length difference between both chromosomes.
func (m imageDNAMater) Mate(parent1, parent2 framework.Candidate,
	numberOfCrossoverPoints int64,
	rng *rand.Rand) []framework.Candidate {

	p1, p2 := parent1.(*imageDNA), parent2.(*imageDNA)

	offspring1 := p1.clone()
	offspring2 := p2.clone()

	var p1min, p1max, p2min, p2max, crossIdx, shorterLen int

	// Apply as many crossovers as required.
	for i := int64(0); i < numberOfCrossoverPoints; i++ {
		p1max = len(offspring1.polys)
		p2max = len(offspring2.polys)
		shorterLen = min(p1max, p2max)
		if p1max == p2max {
			p1min = 0
			p2min = 0
		} else if p1max < p2max {
			// offspring2 is longer than offspring1
			p1min = 0
			p2min = rng.Intn(p2max - p1max)
		} else {
			// offspring1 is longer than offspring2
			p2min = 0
			p1min = rng.Intn(p1max - p2max)
		}

		crossIdx = 1 + rng.Intn(shorterLen-1)
		for j := 0; j < crossIdx; j++ {
			// swap elements of both offsprings
			offspring1.polys[p1min+j], offspring2.polys[p2min+j] = offspring2.polys[p2min+j], offspring1.polys[p1min+j]
		}
	}
	return []framework.Candidate{offspring1, offspring2}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
