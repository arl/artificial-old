package main

import (
	"fmt"
	"image"
	"image/png"
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
	f1, err := os.Open(appConfig.RefImage)
	check(err)
	defer f1.Close()

	img, err := png.Decode(f1)
	check(err)
	bestImg, err := evolveImage(convertToRGBA(img))
	check(err)

	// save best candidate
	f2, err := os.Create("best.png")
	if err != nil {
		panic(err)
	}
	defer f2.Close()
	png.Encode(f2, bestImg)
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

func evolveImage(img *image.RGBA) (image.Image, error) {
	// pseudo random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// chromosome/image factory
	DNAFactory, err := newImageDNAfactory(img.Bounds().Dx(), img.Bounds().Dy())
	if err != nil {
		return nil, err
	}

	// mutation settings
	mutation, err := newImageDNAMutation()
	if err != nil {
		return nil, err
	}

	// define a selection strategy
	selectionStrategy, err := selection.NewTruncationSelection(selection.WithConstantSelectionRatio(0.1))
	if err != nil {
		return nil, err
	}

	// define a fitness evaluator
	evaluator := &fitnessEvaluator{img}

	engine := evolve.NewGenerationalEvolutionEngine(DNAFactory,
		mutation,
		evaluator,
		selectionStrategy,
		rng)

	// define termination conditions
	userAbort := termination.NewUserAbort()
	targetFitness := termination.NewTargetFitness(0, false)

	// define evolution observers
	bestObs, err := newBestObserver(100)
	if err != nil {
		return nil, err
	}
	engine.AddEvolutionObserver(bestObs)


	go func() {
		// handle user termination
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		userAbort.Abort()
	}()

	best := engine.Evolve(
		appConfig.Population.NumIndividuals,
		appConfig.Population.EliteCount,
		userAbort, targetFitness)

	satisfied, err := engine.SatisfiedTerminationConditions()
	if err != nil {
		return nil, err
	}
	fmt.Println("Evolution ended...")
	for _, cond := range satisfied {
		fmt.Println(cond)
	}

	fmt.Println(best)
	return best.(*imageDNA).render(), nil
}

type bestObserver struct {
	frequency int // print statistics every N generations
}

func newBestObserver(frequency int) (o *bestObserver, err error) {
	if frequency == 0 {
		return nil, fmt.Errorf("bessObserver frequency can't be 0")
	}
	return &bestObserver{frequency: frequency}, nil
}

func (o *bestObserver) PopulationUpdate(data *framework.PopulationData) {
	genNum := data.GenerationNumber()
	if genNum%o.frequency == 0 {
		// update best candidate
		fmt.Printf("Generation %d: best: %.2f mean: %.2f stddev: %.2f\n",
			data.GenerationNumber(), data.BestCandidateFitness(), data.MeanFitness(), data.FitnessStandardDeviation())
	}
}
