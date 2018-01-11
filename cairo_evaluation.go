package main

/*
 #cgo pkg-config: cairo
 #include <stdlib.h>
 #include "cairo_evaluation.h"
*/
import "C"
import (
	"image"
	"unsafe"

	"github.com/aurelien-rainone/evolve/framework"
)

type cairoEvaluator struct {
	img *image.RGBA // reference image
}

func (ev *cairoEvaluator) Fitness(cand framework.Candidate, pop []framework.Candidate) float64 {
	var (
		dna    *imageDNA
		cdna   C.imageDNA
		bounds = ev.img.Bounds() // image bounds
	)

	dna = cand.(*imageDNA)

	cdna.w = C.uint32(bounds.Dx())
	cdna.h = C.uint32(bounds.Dy())

	// allocate an array of npolys C structs, of type C.poly
	cdna.npolys = C.uint32(len(dna.polys))
	//fmt.Println("cdna.npolys=", cdna.npolys, "len(dna.polys):", len(dna.polys))
	cdna.polys = (*C.poly)(C.malloc(C.size_t(cdna.npolys) * C.sizeof_poly))
	defer C.free(unsafe.Pointer(cdna.polys))
	cdnaSize := unsafe.Sizeof(cdna)
	var r, g, b, a uint32

	// fill the C.poly's
	for i := uintptr(0); i < uintptr(cdna.npolys); i++ {
		ithPolyAddr := uintptr(unsafe.Pointer(cdna.polys)) + cdnaSize*i
		poly := dna.polys[i]
		cpoly := (*C.poly)(unsafe.Pointer(ithPolyAddr))

		r, g, b, a = poly.col.RGBA()
		cpoly.r = C.uchar(r)
		cpoly.g = C.uchar(g)
		cpoly.b = C.uchar(b)
		cpoly.a = C.uchar(a)
		cpoly.npts = C.uint32(len(poly.pts))

		// allocate an array of npts C struct, of type C.point. A C.point struct
		// being made of 2 int32 points, its size is 8 bytes.

		// TODO: the first element to the slice of points could easily be passed
		// to C without more allocations as the Go and C Point struct layouts
		// are probably the same, or they must be made the same of that is not
		// the case. This would save us a lot of malloc/free

		// fill the C.point's array
		cpoly.pts = (*C.point)(C.malloc(C.size_t(len(poly.pts)) * 8))
		defer C.free(unsafe.Pointer(cpoly.pts))
		for j := uintptr(0); j < uintptr(cpoly.npts); j++ {
			pt := poly.pts[j]
			jthPointAddr := uintptr(unsafe.Pointer(cpoly.pts)) + 8*j
			cpt := (*C.point)(unsafe.Pointer(jthPointAddr))
			cpt.x = C.int32(pt.X)
			cpt.y = C.int32(pt.Y)
		}
	}

	C.render(&cdna)
	return 1
}

func (ev *cairoEvaluator) IsNatural() bool {
	// the lesser the fitness the better
	return false
}

// cgo links:
// https://stackoverflow.com/questions/19910647/pass-struct-and-array-of-structs-to-c-function-from-go
// https://coderwall.com/p/m_ma7q/pass-go-slices-as-c-array-parameters
// cairo doc: file:///usr/share/gtk-doc/html/cairo/cairo-Image-Surfaces.html#cairo-format-t
