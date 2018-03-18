package main

/*
 #cgo pkg-config: cairo
 #include <stdlib.h>
 #include "cairo_evaluation.h"
*/
import "C"
import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"runtime/pprof"
	"time"

	"github.com/aurelien-rainone/evolve"
	"github.com/aurelien-rainone/evolve/operators"
	"github.com/aurelien-rainone/evolve/selection"
	"github.com/aurelien-rainone/evolve/termination"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	err := readConfig()
	check(err)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Info().Msgf("creating cpuprofile %s", *cpuprofile)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Info().Msgf("Reference image %s", appConfig.RefImage)
	f, err := os.Open(appConfig.RefImage)
	check(err)
	defer f.Close()

	img, err := png.Decode(f)
	check(err)
	// start evolution, saving the best candidate as 'best.png'
	err = evolveImage(convertToRGBA(img), "best.png")
	check(err)
}

func saveToPng(fn string, img image.Image) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	png.Encode(f, img)
	return nil
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

func evolveImage(img *image.RGBA, bestPath string) error {
	// pseudo random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// chromosome/image factory
	DNAFactory, err := newImageDNAfactory(img.Bounds().Dx(), img.Bounds().Dy())
	if err != nil {
		return err
	}

	// mutation settings
	mutation, err := newImageDNAMutation()
	if err != nil {
		return err
	}

	// crossover settings
	crossover, err := newImageDNACrossover()
	if err != nil {
		return err
	}

	// create a pipeline that applies mutation then crossover
	pipeline, err := operators.NewEvolutionPipeline(mutation, crossover)
	check(err)

	// define a selection strategy
	selectionStrategy := selection.Identity{}
	//selectionStrategy, err := selection.NewTruncationSelection(selection.WithConstantSelectionRatio(0.1))
	//if err != nil {
	//return nil, err
	//}

	// define a fitness evaluator
	evaluator := newCairoEvaluator(appConfig.RefImage)
	if evaluator == nil {
		return fmt.Errorf("can't create cairo evaluator")
	}

	engine := evolve.NewGenerationalEvolutionEngine(DNAFactory,
		pipeline,
		evaluator,
		selectionStrategy,
		rng)

	engine.SetSingleThreaded(true)

	// define termination conditions
	userAbort := termination.NewUserAbort()
	targetFitness := termination.NewTargetFitness(0, false)

	// define output directory for saves images and generations database
	outDir, err := ioutil.TempDir("_output", "")
	if err != nil {
		return fmt.Errorf("output directory error:, %v", err)
	}
	log.Info().Msgf("ouput directory: %s", outDir)

	// define evolution observers
	bestObs, err := newBestObserver(100, outDir)
	if err != nil {
		return err
	}
	engine.AddEvolutionObserver(bestObs)

	sqliteObs, err := newSqliteObserver(100, outDir)
	if err != nil {
		return err
	}
	defer sqliteObs.close()
	engine.AddEvolutionObserver(sqliteObs)

	// save a copy of refernce image in output dir
	saveToPng(path.Join(outDir, "_ref.png"), img)

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
		return err
	}
	log.Info().Msg("Evolution ended...")
	for _, cond := range satisfied {
		log.Info().Msg(cond.String())
	}

	renderAndDiff(best.(*imageDNA), &bestPath)
	return nil
}
