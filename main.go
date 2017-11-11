package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/aurelien-rainone/evolve"
	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/selection"
	"github.com/aurelien-rainone/evolve/termination"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := readConfig()
	check(err)

	fmt.Println("Reference image:", appConfig.RefImage)
	f, err := os.Open(appConfig.RefImage)
	check(err)
	defer f.Close()

	img, err := png.Decode(f)
	check(err)
	check(evolveImage(convertToRGBA(img)))
}

// convert any image.Image into *image.RGBA
func convertToRGBA(img image.Image) *image.RGBA {
	var rgba *image.RGBA

	switch cimg := img.(type) {
	case *image.RGBA:
		// nothing to do
		rgba = cimg
	default:
		b := img.Bounds()
		rgba = image.NewRGBA(b)
		// convert pixel by pixel
		for y := 0; y < b.Max.Y; y++ {
			for x := 0; x < b.Max.X; x++ {
				col := img.At(x, y)
				rgba.Set(x, y, col)
			}
		}
	}
	return rgba
}

func evolveImage(img *image.RGBA) error {
	// pseudo random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// chromosome/image factory
	DNAFactory, err := newImageDNAfactory(img.Bounds().Dx(), img.Bounds().Dy())
	if err != nil {
		return nil
	}

	// mutation settings
	mutation, err := newImageDNAMutation()
	if err != nil {
		return err
	}

	// define a selection strategy
	var selectionStrategy = &selection.RouletteWheelSelection{}

	// define a fitness evaluator
	evaluator := &fitnessEvaluator{img}

	engine := evolve.NewGenerationalEvolutionEngine(DNAFactory,
		mutation,
		evaluator,
		selectionStrategy,
		rng)

	var obs observer
	engine.AddEvolutionObserver(&obs)

	go func() {
		result := engine.Evolve(2, 1, termination.NewTargetFitness(0, false))
		fmt.Println("Evolution ended...", result)
	}()

	// handle termination
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	log.Println("Evolution interrupted!")

	// save best candidate
	f, err := os.Create("best.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Println(obs.best)
	png.Encode(f, obs.best.render())

	// do last actions and wait for all write operations to end
	os.Exit(0)
	return nil
}

type observer struct {
	best *imageDNA
}

func (o *observer) PopulationUpdate(data *framework.PopulationData) {
	dna := data.BestCandidate().(*imageDNA)
	//o.best = dna.clone()
	o.best = dna
	fmt.Printf("Generation %d: (%v)\n", data.GenerationNumber(), data.BestCandidateFitness())
}
