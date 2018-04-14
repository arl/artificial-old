package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"time"

	"github.com/aurelien-rainone/ARTificial/colors"
	colorful "github.com/lucasb-eyer/go-colorful"
)

func check(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func writePNG(img image.Image, fn string) {
	f, err := os.Create(fn)
	check(err)
	defer f.Close()
	check(png.Encode(f, img))
}

func main() {
	rng := rand.New(rand.NewSource(int64(time.Now().UnixNano())))

	orgc := colorful.MakeColor(
		colors.Lab2RGB[rng.Intn(colors.Resolution)][rng.Intn(colors.Resolution)][rng.Intn(colors.Resolution)])
	dstc := colorful.MakeColor(
		colors.Lab2RGB[rng.Intn(colors.Resolution)][rng.Intn(colors.Resolution)][rng.Intn(colors.Resolution)])

	const (
		nblocks = 10
		blockw  = 80
	)

	img := image.NewRGBA(image.Rect(0, 0, nblocks*blockw, blockw))

	// draw each block wth an increasing degree of blending between
	// the original and destination color
	for i := 0; i < nblocks; i++ {
		l, a, b := orgc.BlendLab(dstc, float64(i)/float64(nblocks-1)).Clamped().Lab()
		lidx := colors.LToIndex(l)
		aidx := colors.AToIndex(a)
		bidx := colors.BToIndex(b)

		draw.Draw(img, image.Rect(i*blockw, 0, (i+1)*blockw, blockw),
			&image.Uniform{colors.Lab2RGB[lidx][aidx][bidx]},
			image.ZP, draw.Src)
	}
	writePNG(img, "colorblend.png")
}
