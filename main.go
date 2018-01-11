package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
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
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := readConfig()
	check(err)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("creating cpuprofile:", *cpuprofile)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Println("Reference image:", appConfig.RefImage)
	f, err := os.Open(appConfig.RefImage)
	check(err)
	defer f.Close()

	img, err := png.Decode(f)
	check(err)
	bestImg, err := evolveImage(convertToRGBA(img))
	check(err)

	// save best candidate
	err = saveToPng("best.png", bestImg)
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

	// crossover settings
	crossover, err := newImageDNACrossover()
	if err != nil {
		return nil, err
	}

	// create a pipeline that applies mutation then crossover
	pipeline, err := operators.NewEvolutionPipeline(mutation, crossover)
	check(err)

	// define a selection strategy
	selectionStrategy := selection.Identity{}

	// define a fitness evaluator
	evaluator := &cairoEvaluator{img}

	engine := evolve.NewGenerationalEvolutionEngine(DNAFactory,
		pipeline,
		evaluator,
		selectionStrategy,
		rng)

	// define termination conditions
	userAbort := termination.NewUserAbort()
	targetFitness := termination.NewTargetFitness(0, false)

	// define output directory for saves images and generations database
	outDir, err := ioutil.TempDir("_output", "")
	if err != nil {
		return nil, fmt.Errorf("output directory error:, %v", err)
	}
	log.Println("ouput directory:", outDir)

	// define evolution observers
	bestObs, err := newBestObserver(100, outDir)
	if err != nil {
		return nil, err
	}
	engine.AddEvolutionObserver(bestObs)

	sqliteObs, err := newSqliteObserver(100, outDir)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	fmt.Println("Evolution ended...")
	for _, cond := range satisfied {
		fmt.Println(cond)
	}
	return best.(*imageDNA).render(), nil
}
