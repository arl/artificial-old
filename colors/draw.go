package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	cimg := image.NewRGBA(image.Rect(0, 0, 64, 64))

	halfpoint := Lab2RGBNumSteps / 4

	for x := 0; x < Lab2RGBNumSteps; x++ {
		for y := 0; y < Lab2RGBNumSteps; y++ {
			cimg.SetRGBA(x, y, Lab2RGB[x][y][halfpoint])
		}
	}

	f, err := os.Create(fmt.Sprintf("palette_%v_%v.png", halfpoint, Lab2RGBNumSteps))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, cimg)

	if err != nil {
		log.Fatal(err)
	}
}
