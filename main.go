package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
)

var inputFile string

func init() {
	flag.StringVar(&inputFile, "input", "", "reference image")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	if len(inputFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Println("loading inputFile")
	infile, err := os.Open(inputFile)
	check(err)
	defer infile.Close()

	img, err := png.Decode(infile)
	check(err)

	check(evolve(img))
}

func evolve(img image.Image) error {
	// gene factory
	DNAFactory, err := newImageDNAfactory(50, img.Bounds().Dx(), img.Bounds().Dy())
	if err != nil {
		return nil
	}
	fmt.Println(DNAFactory)
	return nil
}
