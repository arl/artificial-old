package main

/*
 #cgo pkg-config: cairo
 #include <stdlib.h>
 #include "cairo_evaluation.h"
*/
import "C"
import (
	"unsafe"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/rs/zerolog/log"
)

type cairoEvaluator struct {
	orgImgW, orgImgH C.uint32
}

func newCairoEvaluator(path string) *cairoEvaluator {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ev := new(cairoEvaluator)
	rc := C.evaluator_init(cpath, 128, 128, &ev.orgImgW, &ev.orgImgH)
	C.fflush(C.stdout)
	if rc != 0 {
		log.Fatal().Msgf("evaluator_init returned %v\n", rc)
		return nil
	}
	return ev
}

func (ev *cairoEvaluator) Fitness(cand framework.Candidate, pop []framework.Candidate) float64 {
	var (
		dna  *imageDNA
		cdna C.imageDNA
	)

	dna = cand.(*imageDNA)

	cdna.w = ev.orgImgW
	cdna.h = ev.orgImgH

	// allocate an array of npolys C structs, of type C.poly
	cdna.npolys = C.uint32(len(dna.polys))
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

	var diffval C.double = 0
	rc := C.render_and_diff(&cdna, &diffval, 0, false)
	if rc != 1 {
		log.Fatal().Msg("render_and_diff errored")
	}
	return float64(diffval)
}

func (ev *cairoEvaluator) IsNatural() bool {
	// the lesser the fitness the better
	return false
}
