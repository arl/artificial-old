package main

import (
	"flag"
	"fmt"

	cfg "github.com/jinzhu/configor"
)

var appConfig = struct {
	// path to the reference image
	RefImage string `required:"true"`

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
	}
}{}

func readConfig() error {
	configFile := flag.String("cfg", "config.yml", "configuration file")
	flag.Parse()
	err := cfg.Load(&appConfig, *configFile)
	if err != nil {
		return fmt.Errorf("read config error: %v", err)
	}
	fmt.Println(appConfig)
	return nil
}
