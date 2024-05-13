package earthly

import (
	"bytes"
	"image"
	"image/png"
	"math"
)

type EarthlyConfig struct {
	Size       int
	Background []uint8
}

func (config *EarthlyConfig) Generate() *bytes.Buffer {
	canvas := image.NewNRGBA(image.Rect(0, 0, config.Size, config.Size))

	var halfSize float64 = float64(config.Size) / 2.0
	var onePxSize = 1.0 / halfSize
	var twoPxSize = onePxSize * 2.0

	// A circle of 1.0 radius would touch the very edges of the image
	// In some cases we want to have border effects on the image so we make it a little smaller
	var circleRadius = 0.9

	for py := 0; py < config.Size; py++ {
		for px := 0; px < config.Size; px++ {
			var r, g, b, a uint8 = config.Background[0], config.Background[1], config.Background[2], config.Background[3]

			var u, v float64 = (float64(px) - halfSize) * onePxSize, -(float64(py) - halfSize) * onePxSize
			magnitude := math.Sqrt(math.Pow(u, 2.0) + math.Pow(v, 2.0))
			circleMagnitude := magnitude / circleRadius
			var cx, cy = u / circleRadius, v / circleRadius
			var cz = 0.0

			if circleMagnitude < 1.0 {
				cz = math.Sqrt(1.0 - math.Pow(circleMagnitude, 2.0))
			} else {
				cx = cx / circleMagnitude
				cy = cy / circleMagnitude
			}

			longitude := math.Asin(cy)
			latitude := math.Atan2(-cx, cz)

			if math.Mod(((longitude*(180/math.Pi))+360), 30) > 15 {
				r = 128
			}

			if math.Mod(((latitude*(180/math.Pi))+360), 30) > 15 {
				g = 128
			}

			// Mask it into a soft circle shape
			circleMask := 1.0 - math.Max(0.0, math.Min(1.0, (circleMagnitude-(1.0-twoPxSize))/twoPxSize))

			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+0] = r
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+1] = g
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+2] = b
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+3] = uint8(float64(a) * circleMask)
		}
	}

	buffer := new(bytes.Buffer)
	png.Encode(buffer, canvas)

	return buffer
}
