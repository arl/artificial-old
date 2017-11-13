package main

import (
	"flag"
	"fmt"

	cfg "github.com/jinzhu/configor"
)

var appConfig = struct {
	// path to the reference image
	RefImage string

	Population struct {
		// number of individuals in the population
		NumIndividuals int `required:"true"`
		// number of candidates preserved through elitism
		EliteCount int `required:"true"`
	}

	Image struct {
		// MinPolys is the minimum number of polygon in an image
		MinPolys int `required:"true"`
		// MaxPolys is the minimum number of polygon in an image
		MaxPolys int `required:"true"`
	}

	Polygon struct {
		// MinPoints is the minimum number of points in a polygon
		MinPoints int `required:"true"`
		// MaxPoints is the maximum number of points in a polygon
		MaxPoints int `required:"true"`
	}

	Mutation struct {
		// image level mutations
		Image struct {
			// Rate [0, 1] of add polygon mutation
			AddPoly float64 `required:"true"`
			// Rate [0, 1] of remove polygon mutation
			RemovePoly float64 `required:"true"`
			// Rate [0, 1] of swap polygon mutation
			SwapPolys float64 `required:"true"`
			// Rate [0, 1] of background color mutation
			BackgroundColor float64 `required:"true"`
		}

		// polygon level mutations
		Polygon struct {
			// Rate [0, 1] of add point mutation
			AddPoint float64 `required:"true"`
			// Rate [0, 1] of remove point mutation
			RemovePoint float64 `required:"true"`
			// Rate [0, 1] of change polygon color mutation
			ChangeColor float64 `required:"true"`
		}

		// point level mutations
		Point struct {
			// Rate [0, 1] of move point mutation
			Move float64 `required:"true"`
		}
	}
}{}

func readConfig() error {
	configFile := flag.String("cfg", "config.yml", "configuration file")
	refImage := flag.String("img", "", "reference image (PNG)")
	flag.Parse()
	err := cfg.Load(&appConfig, *configFile)
	if err != nil {
		return fmt.Errorf("read config error: %v", err)
	}
	fmt.Println(appConfig)
	if len(*refImage) > 0 {
		appConfig.RefImage = *refImage
	}
	return nil
}
