package earthly

import (
	"bytes"
	"image"
	"image/png"
)

type EarthlyConfig struct {
	Size		int
	Background	[]uint8
}

func (config *EarthlyConfig) Generate() (*bytes.Buffer) {
	canvas := image.NewRGBA(image.Rect(0, 0, config.Size, config.Size))

	for y := 0; y < config.Size; y++ {
		for x := 0; x < config.Size; x++ {

			canvas.Pix[(y-canvas.Rect.Min.Y)*canvas.Stride + (x-canvas.Rect.Min.X)*4 + 0] = config.Background[0];
			canvas.Pix[(y-canvas.Rect.Min.Y)*canvas.Stride + (x-canvas.Rect.Min.X)*4 + 1] = config.Background[1];
			canvas.Pix[(y-canvas.Rect.Min.Y)*canvas.Stride + (x-canvas.Rect.Min.X)*4 + 2] = config.Background[2];
			canvas.Pix[(y-canvas.Rect.Min.Y)*canvas.Stride + (x-canvas.Rect.Min.X)*4 + 3] = config.Background[3];

		}
	}

	buffer := new(bytes.Buffer)
	png.Encode(buffer, canvas)

	return buffer
}
