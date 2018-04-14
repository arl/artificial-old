package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/aurelien-rainone/ARTificial/colors"
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
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	halfpoint := colors.Resolution / 2

	for x := 0; x < colors.Resolution; x++ {
		for y := 0; y < colors.Resolution; y++ {
			img.SetRGBA(x, y, colors.Lab2RGB[x][y][halfpoint])
		}
	}

	writePNG(img, fmt.Sprintf("palette_%v_%v.png", halfpoint, colors.Resolution))
}
