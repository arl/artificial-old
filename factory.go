package main

import (
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve/factory"
	"github.com/aurelien-rainone/evolve/framework"
)

type imageDNAfactory struct {
	factory.AbstractCandidateFactory
}

func newImageDNAfactory(numPolys, imgW, imgH int) (*imageDNAfactory, error) {
	if imgW == 0 || imgH == 0 {
		return nil, fmt.Errorf("invalid dimensions %v x %v", imgW, imgH)
	}

	sf := &imageDNAfactory{
		factory.AbstractCandidateFactory{
			&imageDNAGenerator{
				numPolys: numPolys,
				imgW:     imgW,
				imgH:     imgH,
			},
		},
	}
	return sf, nil
}

type imageDNAGenerator struct {
	numPolys   int
	imgW, imgH int
}

func (g *imageDNAGenerator) GenerateRandomCandidate(rng *rand.Rand) framework.Candidate {
	var img imageDNA
	img.polys = make([]poly, g.numPolys)
	for i := 0; i < g.numPolys; i++ {
		img.polys[i] = randomPoly(&img, newPolyMinPoints, newPolyMaxPoints, rng)
	}
	return &img
}
