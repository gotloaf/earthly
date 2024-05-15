package earthly

import (
	"bytes"
	"image"
	"image/png"
	"math"
)

type EarthlyConfig struct {
	Size       	int
	Background 	[]uint8
	Latitude   	float64
	Longitude  	float64
	Roll	   	float64
}

func (config *EarthlyConfig) Generate() *bytes.Buffer {
	canvas := image.NewNRGBA(image.Rect(0, 0, config.Size, config.Size))

	var halfSize float64 = float64(config.Size) / 2.0
	var onePxSize = 1.0 / halfSize
	var twoPxSize = onePxSize * 2.0

	// A circle of 1.0 radius would touch the very edges of the image
	// In some cases we want to have border effects on the image so we make it a little smaller
	var circleRadius = 0.9

	// Pre-calculate the longitude and latitude's rotation matrix coefficients
	latitudeRad := config.Latitude * (math.Pi / 180)
	latitudeSin := math.Sin(latitudeRad)
	latitudeCos := math.Cos(latitudeRad)
	longitudeRad := config.Longitude * (math.Pi / 180)
	longitudeSin := math.Sin(longitudeRad)
	longitudeCos := math.Cos(longitudeRad)
	rollRad := config.Roll * (math.Pi / 180)
	rollSin := math.Sin(rollRad)
	rollCos := math.Cos(rollRad)

	matrix := []float64{
		(latitudeCos * rollCos + latitudeSin * longitudeSin * rollSin), (latitudeSin * longitudeSin * rollCos - latitudeCos * rollSin), -(latitudeSin * longitudeCos),
		(longitudeCos * rollSin), (longitudeCos * rollCos), (longitudeSin),
		(latitudeSin * rollCos - latitudeCos * longitudeSin * rollSin), (-latitudeCos * longitudeSin * rollCos - latitudeSin * rollSin), (latitudeCos * longitudeCos),
	}


	for py := 0; py < config.Size; py++ {
		for px := 0; px < config.Size; px++ {
			var r, g, b uint8 = 0, 0, 0
			var a float64 = 0.0

			var u, v float64 = (float64(px) - halfSize) * onePxSize, -(float64(py) - halfSize) * onePxSize
			magnitude := math.Sqrt(math.Pow(u, 2.0) + math.Pow(v, 2.0))
			circleMagnitude := magnitude / circleRadius
			var sphereX1, sphereY1 = u / circleRadius, v / circleRadius
			var sphereZ1 = 0.0

			if circleMagnitude < 1.0 {
				sphereZ1 = math.Sqrt(1.0 - math.Pow(circleMagnitude, 2.0))
			} else {
				sphereX1 = sphereX1 / circleMagnitude
				sphereY1 = sphereY1 / circleMagnitude
			}

			sphereX2 := matrix[0] * sphereX1 + matrix[1] * sphereY1 + matrix[2] * sphereZ1
			sphereY2 := matrix[3] * sphereX1 + matrix[4] * sphereY1 + matrix[5] * sphereZ1
			sphereZ2 := matrix[6] * sphereX1 + matrix[7] * sphereY1 + matrix[8] * sphereZ1

			projectedLongitude := math.Asin(sphereY2)
			projectedLatitude := math.Atan2(-sphereX2, sphereZ2)

			if math.Mod(((projectedLongitude*(180/math.Pi))+360), 30) > 15 {
				r = 128
			}

			if math.Mod(((projectedLatitude*(180/math.Pi))+360), 30) > 15 {
				g = 128
			}

			// Mask it into a soft circle shape
			circleMask := 1.0 - math.Max(0.0, math.Min(1.0, (circleMagnitude-(1.0-twoPxSize))/twoPxSize))

			a = circleMask

			// Composite onto background
			bg_a := float64(config.Background[3]) / 255.0
			bg_r := float64(config.Background[0]) * bg_a
			bg_g := float64(config.Background[1]) * bg_a
			bg_b := float64(config.Background[2]) * bg_a

			final_alpha := bg_a + a - bg_a*a
			final_r := (float64(r)*a + bg_r*(1.0-a)) / final_alpha
			final_g := (float64(g)*a + bg_g*(1.0-a)) / final_alpha
			final_b := (float64(b)*a + bg_b*(1.0-a)) / final_alpha

			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+0] = uint8(final_r)
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+1] = uint8(final_g)
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+2] = uint8(final_b)
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+3] = uint8(final_alpha * 255.0)
		}
	}

	buffer := new(bytes.Buffer)
	png.Encode(buffer, canvas)

	return buffer
}
