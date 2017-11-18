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

func newImageDNAfactory(imgW, imgH int) (*imageDNAfactory, error) {
	if imgW == 0 || imgH == 0 {
		return nil, fmt.Errorf("invalid dimensions %v x %v", imgW, imgH)
	}

	sf := &imageDNAfactory{
		factory.AbstractCandidateFactory{
			RandomCandidateGenerator: &imageDNAGenerator{
				imgW: imgW,
				imgH: imgH,
			},
		},
	}
	return sf, nil
}

type imageDNAGenerator struct {
	imgW, imgH int // width/height of the reference image
}

func (g *imageDNAGenerator) GenerateRandomCandidate(rng *rand.Rand) framework.Candidate {
	var numPolys int
	if appConfig.Image.MinPolys == appConfig.Image.MaxPolys {
		numPolys = appConfig.Image.MaxPolys
	} else {
		numPolys = appConfig.Image.MinPolys + rng.Intn(appConfig.Image.MaxPolys-appConfig.Image.MinPolys)
	}

	// create image dna with same dimensions than reference image
	var img = &imageDNA{
		w:     g.imgW,
		h:     g.imgH,
		polys: make([]poly, numPolys),
	}
	// add N `numPolys` random polygons
	for i := 0; i < numPolys; i++ {
		img.polys[i] = randomPoly(img, appConfig.Polygon.MinPoints, appConfig.Polygon.MaxPoints, rng)
	}
	return img
}
