package main

import (
	"fmt"
	"image/color"
	"io"

	colorful "github.com/lucasb-eyer/go-colorful"
)

const res = 256

var lut [res][res][res][3]byte

func main3() {
	c := colorful.Lab(0.3, 0.3, 0.3)
	r, g, b := c.RGB255()
	labL, labA, labB := colorful.MakeColor(color.RGBA{r, g, b, 255}).Lab()
	fmt.Println(labL, labA, labB)
}

func writeLUT(w io.Writer) error {
	// for each color component, the number of points describing the full range.
	const nsteps = 2
	var lut [nsteps][nsteps][nsteps][3]byte
	lstep := (1.0 - 0.0) / (nsteps - 1)
	abstep := (1.0 - -1.0) / (nsteps - 1)

	//f.Write([]byte("var labToRGBLut "))

	for lidx := 0; lidx < nsteps; lidx++ {
		l := float64(lidx) * lstep
		for aidx := 0; aidx < nsteps; aidx++ {
			a := float64(aidx)*abstep - 1.0
			for bidx := 0; bidx < nsteps; bidx++ {
				b := float64(bidx)*abstep - 1.0
				R, G, B := colorful.Lab(l, a, b).RGB255()
				lut[lidx][aidx][bidx] = [3]byte{R, G, B}
				//w.Write()
			}
		}
	}
	fmt.Println(lut)

	//RGBA := lut[0][0][0]
	//col := color.RGBA{RGBA[0], RGBA[1], RGBA[2], 255}
	//c := colorful.MakeColor(col)
	//fmt.Println(RGBA)
	//fmt.Println(c.RGBA())
	return nil
}

func main() {

	const nsteps = 2
	lut := [nsteps][nsteps][nsteps]color.Color{
		{
			{
				color.RGBA{},
				color.RGBA{},
			},
			{
				color.RGBA{},
				color.RGBA{},
			},
		},
		{
			{
				color.RGBA{},
				color.RGBA{},
			},
			{
				color.RGBA{},
				color.RGBA{},
			},
		},
	}
	fmt.Printf("%#v", lut)

	//err := writeLUT(os.Stdout)
	//if err != nil {
	//fmt.Println(err)
	//}
}

func main2() {
	//c := colorful.Lab(0.507850, 0.040585, -0.370945)

	// CIE-L*a*b*: A perceptually uniform color space, i.e. distances are
	// meaningful. L* in [0..1] and a*, b* almost in [-1..1].

	lstep := 1.0 / res
	abstep := 2.0 / res
	for R, l := 0, 0.0; l < 1.0; R, l = R+1, l+lstep {
		//fmt.Println("l", l)
		for G, a := 0, -1.0; a < 1.0; G, a = G+1, a+abstep {
			for B, b := 0, -1.0; b < 1.0; B, b = B+1, b+abstep {
				c := colorful.Lab(l, a, b)
				r_, g_, b_ := c.RGB255()
				lut[R][G][B][0] = r_
				lut[R][G][B][1] = g_
				lut[R][G][B][2] = b_
			}
		}
	}

	// test values
	tl, ta, tb := 0.25, 0.10, 0.70

	// look for the corresponding RGB in the lookup table
	cRGB := lut[int(tl*lstep)][int(ta*abstep)][int(tb*abstep)]

	fmt.Println("from lut    : ", cRGB)
	R, G, B := colorful.Lab(tl, ta, tb).RGB255()
	fmt.Println("from library: ", R, G, B)

	l, a, b := colorful.MakeColor(color.RGBA{R, G, B, 255}).Lab()
	fmt.Println("and back: ", l, a, b)

	//fmt.Println("from lut => ", 100*lstep, 25*abstep, 56*abstep)
	//fmt.Println("real lab => ", l, a, b)
	//toimg, err := os.Create("colorblend.png")
	//if err != nil {
	//fmt.Printf("Error: %v", err)
	//return
	//}
	//defer toimg.Close()

	//png.Encode(toimg, img)
}
