package earthly

import (
	"bytes"
	"image"
	"image/png"
)

type EarthlyConfig struct {
	fill	string
}

func (config *EarthlyConfig) Generate() (*bytes.Buffer) {

	canvas := image.NewRGBA(image.Rect(0, 0, 512, 512))

	buffer := new(bytes.Buffer)
	png.Encode(buffer, canvas)

	return buffer
}
