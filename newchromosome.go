package main

import (
	"fmt"
	"image"
	"math"

	"github.com/aurelien-rainone/ARTificial/colors"
	"github.com/aurelien-rainone/evolve/bitstring"
)

// XXX: Documentation
// https://developer.apple.com/library/content/documentation/GraphicsImaging/Conceptual/drawingwithquartz2d/dq_images/dq_images.html

// Extracting phenotype (polygons definition) from genotype (the bitstring)
// see requiredBits

// TODO: those are constants for now but should really be in config
const (
	totalPolygons = 50
)

// XXX: those are constants and should remain as such
const (
	// polygons are at most 7 sided
	maxSegments = 7
)

var (
	bpc      int // bits per color component
	workingw int // working width
	workingh int // working height
)

// return the number 32 bits integers required, with the current configuration,
// to define an image in term of a bistring
// TODO: should take into account current config...
func requiredInt32s() (int, error) {
	var (
		req int
		// TODO: should place it in the globals var right? to make it accessible
		// to the rendering function?
		bpd int // number of bits per dimension
		ok  bool
	)

	if workingw != workingh {
		return 0, fmt.Errorf("working image must be square")
	}
	bpd, ok = ispowerof2(workingw)
	if !ok {
		return 0, fmt.Errorf("working width must be a power of two")
	}

	bpc, ok = ispowerof2(colors.Resolution)
	if !ok {
		return 0, fmt.Errorf("colors.Resolution must be a power of two")
	}
	// 32 bits polygon header
	// 1 bit, bit 0			=> polygon visibility
	// 2 bits, bits 1 to 2	=> nsegments - 3
	//						   the polygon is at least a triangle and at most a
	//						   7 segments polygon
	//
	// That lets us more than enough bits to represent the polygon color.
	// Reasonable bits per component might range from 4 to 8, depending on the
	// desired color range and thus the color precision.
	// Let's accept bpc values in the range [4, 8], that should cover many, if
	// not all, `use cases`.
	if bpc < 4 || bpc > 8 {
		return 0, fmt.Errorf("number of bits ber component must be in [4, 8], got %v", bpc)
	}
	// Required number of bits per polygon point: 32
	// 1 point should be packed on 32 bits, that means:
	// from 1 bit to 16 bits per dimension, that more than enough!
	// For simplicity, x will be coded on the first 16 bits, Y on the next 16
	// bits.
	if bpd < 1 || bpd > 16 {
		return 0, fmt.Errorf("number of bits ber dimension must be in [1, 16], got %v", bpd)
	}
	const ptgeo = 1 // one 32 bits integer

	req = totalPolygons * (1 + maxSegments*ptgeo)

	return req, nil
}

// checks if a is a power of two. If that is the case, the log2 of a is returned
// with true, anf false otherwise.
func ispowerof2(pow int) (n int, ok bool) {
	n = int(math.Logb(float64(pow)))
	fmt.Println(pow, n)
	if 1<<uint(n) == pow {
		return n, true
	}
	return
}

func renderBitstring(bs bitstring.BitString) (image.Image, error) {
	//bits := bs.Data()

	//continuer ici

	return nil
}
