package earthly

import (
	"os"
	"testing"
)

func TestGenerateFromJPEG(t *testing.T) {
	config := EarthlyConfig{
		Size:       1024,
		Background: []uint8{0, 0, 0, 0},
		Latitude:   0.0,
		Longitude:  0.0,
		Roll:       0.0,
		Halo:       true,
		Radius:     1.0,
	}

	image_file, err := os.ReadFile("build/equirectangular/earth_1x.jpg")
	if err != nil {
		t.Errorf("error occurred while reading asset dependency: %s", err)
	}

	output := config.Generate(EarthlyBuffers{
		Earth: image_file,
	}, true)

	if output == nil || output.Len() == 0 {
		t.Errorf("internal error occurred during generation")
	}
}

func TestGenerateFromPNG(t *testing.T) {
	config := EarthlyConfig{
		Size:       1024,
		Background: []uint8{0, 0, 0, 0},
		Latitude:   0.0,
		Longitude:  0.0,
		Roll:       0.0,
		Halo:       true,
		Radius:     1.0,
	}

	image_file, err := os.ReadFile("build/equirectangular/earth_1x.png")
	if err != nil {
		t.Errorf("error occurred while reading asset dependency: %s", err)
	}

	output := config.Generate(EarthlyBuffers{
		Earth: image_file,
	}, false)

	if output == nil || output.Len() == 0 {
		t.Errorf("internal error occurred during generation")
	}
}

func TestGenerateZoomedIn(t *testing.T) {
	config := EarthlyConfig{
		Size:       1024,
		Background: []uint8{0, 0, 0, 0},
		Latitude:   0.0,
		Longitude:  0.0,
		Roll:       0.0,
		Halo:       true,
		Radius:     2.0,
	}

	image_file, err := os.ReadFile("build/equirectangular/earth_1x.png")
	if err != nil {
		t.Errorf("error occurred while reading asset dependency: %s", err)
	}

	output := config.Generate(EarthlyBuffers{
		Earth: image_file,
	}, false)

	if output == nil || output.Len() == 0 {
		t.Errorf("internal error occurred during generation")
	}
}

func TestGenerateZoomedOut(t *testing.T) {
	config := EarthlyConfig{
		Size:       1024,
		Background: []uint8{0, 0, 0, 0},
		Latitude:   0.0,
		Longitude:  0.0,
		Roll:       0.0,
		Halo:       true,
		Radius:     0.5,
	}

	image_file, err := os.ReadFile("build/equirectangular/earth_1x.png")
	if err != nil {
		t.Errorf("error occurred while reading asset dependency: %s", err)
	}

	output := config.Generate(EarthlyBuffers{
		Earth: image_file,
	}, false)

	if output == nil || output.Len() == 0 {
		t.Errorf("internal error occurred during generation")
	}
}

func TestBufferEmptyOnInvalidInput(t *testing.T) {
	config := EarthlyConfig{
		Size:       1024,
		Background: []uint8{0, 0, 0, 0},
		Latitude:   0.0,
		Longitude:  0.0,
		Roll:       0.0,
		Halo:       true,
		Radius:     1.0,
	}

	output := config.Generate(EarthlyBuffers{
		Earth: make([]byte, 0),
	}, false)

	if output != nil && output.Len() > 0 {
		t.Errorf("was unexpectedly able to get output from invalid input")
	}
}
