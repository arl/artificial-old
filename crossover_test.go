package main

import (
	"image/color"
	"math/rand"
	"testing"
)

func TestCrossover(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	img1 := &imageDNA{
		polys: []poly{
			poly{col: color.RGBA{A: 0}},
			poly{col: color.RGBA{A: 1}},
			poly{col: color.RGBA{A: 2}},
			poly{col: color.RGBA{A: 3}},
		},
	}

	img2 := &imageDNA{
		polys: []poly{
			poly{col: color.RGBA{A: 4}},
			poly{col: color.RGBA{A: 5}},
		},
	}
	mater := imageDNAMater{}
	result := mater.Mate(img1, img2, 1, rng)

	offspring1 := result[0].(*imageDNA)
	offspring2 := result[1].(*imageDNA)

	if offspring1.polys[0].col.(color.RGBA).A != 0 {
		t.Errorf("want offspring1.polys[0].col.A = 0")
	}
	if offspring1.polys[1].col.(color.RGBA).A != 4 {
		t.Errorf("want offspring1.polys[1].col.A = 4")
	}
	if offspring1.polys[2].col.(color.RGBA).A != 2 {
		t.Errorf("want offspring1.polys[2].col.A = 2")
	}
	if offspring1.polys[3].col.(color.RGBA).A != 3 {
		t.Errorf("want offspring1.polys[3].col.A = 3")
	}
	if offspring2.polys[0].col.(color.RGBA).A != 1 {
		t.Errorf("want offspring2.polys[0].col.A = 1")
	}
	if offspring2.polys[1].col.(color.RGBA).A != 5 {
		t.Errorf("want offspring2.polys[1].col.A = 5")
	}
}
