package main

import (
	"fmt"
	"image"
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
		numPoints := 3 + rng.Intn(3)
		p := &img.polys[i]
		p.pts = make([]image.Point, numPoints)
		p.col.R = byte(rng.Intn(255))
		p.col.G = byte(rng.Intn(255))
		p.col.B = byte(rng.Intn(255))
		p.col.A = byte(10 + rng.Intn(50))
	}
	return &img
}
