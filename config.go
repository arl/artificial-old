package main

import (
	"flag"
	"fmt"

	cfg "github.com/jinzhu/configor"
)

var appConfig = struct {
	// path to the reference image
	RefImage string

	Image struct {
		// MinPolys is the minimum number of polygon in an image
		MinPolys int `required:"true"`
		// MaxPolys is the minimum number of polygon in an image
		MaxPolys int `required:"true"`
	}

	Polygon struct {
		// MinPoints is the minimum number of points in a polygon
		MinPoints int `required:"true"`
		// MaxPoints is the minimum number of points in a polygon
		MaxPoints int `required:"true"`
	}

	Mutation struct {
		Polygon struct {
			// Rate [0, 1] of add polygon mutation
			Add float64 `required:"true"`
			// Rate [0, 1] of remove polygon mutation
			Remove float64 `required:"true"`
			// Rate [0, 1] of swap polygon mutation
			Swap float64 `required:"true"`
			// Rate [0, 1] of change polygon color mutation
			Color float64 `required:"true"`
		}
		Point struct {
			// Rate [0, 1] of add point mutation
			Add float64 `required:"true"`
			// Rate [0, 1] of remove point mutation
			Remove float64 `required:"true"`
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
