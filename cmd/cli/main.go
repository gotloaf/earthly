package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/swaggo/swag"

	"github.com/gotloaf/earthly"
)

var (
	outputFlag    string
	sizeFlag      int
	longitudeFlag float64
	latitudeFlag  float64
	rollFlag      float64
	zoomFlag      float64
)

func main() {
	flag.StringVar(&outputFlag, "output", "-", "Location to output the image to (defaults to stdout)")
	flag.IntVar(&sizeFlag, "size", 1024, "Size of output image")
	flag.Float64Var(&longitudeFlag, "longitude", 0.0, "Longitude to display")
	flag.Float64Var(&latitudeFlag, "latitude", 0.0, "Latitude to display")
	flag.Float64Var(&rollFlag, "roll", 0.0, "Rotation of display camera")
	flag.Float64Var(&zoomFlag, "zoom", 1.0, "Zoom factor of camera")
	flag.Parse()

	config := earthly.EarthlyConfig{
		Size:       sizeFlag,
		Background: []uint8{0, 0, 0, 0},
		Latitude:   latitudeFlag,
		Longitude:  longitudeFlag,
		Roll:       rollFlag,
		Halo:       true,
		Radius:     zoomFlag,
	}

	var image_format = "1x.png"
	if float64(config.Size)*config.Radius > 1024 {
		image_format = "2x.png"
	}

	image_file, err := os.ReadFile(fmt.Sprintf("build/equirectangular/earth_%s", image_format))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occurred while reading asset dependency: %s", err)
		os.Exit(2) // ENOENT
		return
	}

	output := config.Generate(earthly.EarthlyBuffers{
		Earth: image_file,
	}, false)

	if output == nil || output.Len() == 0 {
		fmt.Fprintf(os.Stderr, "internal error occurred during generation")
		os.Exit(131) // ENOTRECOVERABLE
		return
	}

	if outputFlag == "-" || outputFlag == "" {
		_, err := os.Stdout.Write(output.Bytes())

		if err != nil {
			fmt.Fprintf(os.Stderr, "error occurred while outputting to stdout: %s", err)
			os.Exit(5) // EIO
			return
		}
	} else {
		err := os.WriteFile(outputFlag, output.Bytes(), 0644)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error occurred while outputting file: %s", err)
			os.Exit(5) // EIO
			return
		}
	}
}
