package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"

	"github.com/aurelien-rainone/evolve"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/aurelien-rainone/evolve/selection"
	"github.com/aurelien-rainone/evolve/termination"
)

var (
	inputFile        string
	newPolyMaxPoints int
	newPolyMinPoints int
)

func init() {
	flag.StringVar(&inputFile, "input", "", "reference image (only PNG)")
	flag.IntVar(&newPolyMinPoints, "min-points", 3, "minimum number of points for new polygons")
	flag.IntVar(&newPolyMaxPoints, "max-points", 6, "maximum number of points for new polygons")
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
	check(evolveImage(img))
}

func evolveImage(img image.Image) error {
	// chromosome/image factory
	DNAFactory, err := newImageDNAfactory(50, img.Bounds().Dx(), img.Bounds().Dy())
	if err != nil {
		return nil
	}
	fmt.Println(DNAFactory)

	// mutation
	mutationOptions := mutationOptions{}
	mutationOptions.addPolygonMutation = number.NewConstantProbabilityGenerator(0.1)
	mutationOptions.removePolygonMutation = number.NewConstantProbabilityGenerator(0.1)
	mutationOptions.swapPolygonsMutation = number.NewConstantProbabilityGenerator(0.1)
	mutation, err := newImageDNAMutation(mutationOptions)
	if err != nil {
		return err
	}

	// pseudo random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// define a selection strategy
	var selectionStrategy = &selection.RouletteWheelSelection{}

	// define a fitness evaluator
	evaluator := &fitnessEvaluator{img}

	engine := evolve.NewGenerationalEvolutionEngine(DNAFactory,
		mutation,
		evaluator,
		selectionStrategy,
		rng)

	result := engine.Evolve(10, 5, termination.NewTargetFitness(0, false))
	fmt.Println(result)
	return nil
}
