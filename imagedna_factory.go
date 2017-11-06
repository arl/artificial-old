package main

import (
	"fmt"
	"image"
	"math/rand"

	"github.com/aurelien-rainone/evolve/factory"
	"github.com/aurelien-rainone/evolve/framework"
)

type ImageDNAFactory struct {
	factory.AbstractCandidateFactory
}

func NewImageDNAFactory(numPolys, imgW, imgH int) (*ImageDNAFactory, error) {
	if imgW == 0 || imgH == 0 {
		return nil, fmt.Errorf("invalid dimensions %v x %v", imgW, imgH)
	}

	sf := &ImageDNAFactory{
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
	var img ImageDNA
	img.polys = make([]Poly, g.numPolys)
	for i := 0; i < g.numPolys; i++ {
		numPoints := 3 + rng.Intn(3)
		poly := &img.polys[i]
		poly.pts = make([]image.Point, numPoints)
		poly.col.R = byte(rng.Intn(255))
		poly.col.G = byte(rng.Intn(255))
		poly.col.B = byte(rng.Intn(255))
		poly.col.A = byte(10 + rng.Intn(50))
	}
	return &img
}
