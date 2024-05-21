package earthly

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"math"
)

type EarthlyConfig struct {
	Size       int
	Background []uint8
	Latitude   float64
	Longitude  float64
	Roll       float64
	Halo       bool
	Radius     float64
}

type EarthlyBuffers struct {
	Earth []byte
}

func AlphaComposite(
	fg_r uint8, fg_g uint8, fg_b uint8, fg_a uint8,
	bg_r uint8, bg_g uint8, bg_b uint8, bg_a uint8,
) (r uint8, g uint8, b uint8, a uint8) {

	bg_af := float64(bg_a) / 255.0
	bg_rf := float64(bg_r) * bg_af
	bg_gf := float64(bg_g) * bg_af
	bg_bf := float64(bg_b) * bg_af

	af := float64(fg_a) / 255.0

	final_alpha := bg_af + af - bg_af*af
	final_r := (float64(fg_r)*af + bg_rf*(1.0-af)) / final_alpha
	final_g := (float64(fg_g)*af + bg_gf*(1.0-af)) / final_alpha
	final_b := (float64(fg_b)*af + bg_bf*(1.0-af)) / final_alpha

	return uint8(final_r), uint8(final_g), uint8(final_b), uint8(final_alpha * 255.0)
}

func (config *EarthlyConfig) Generate(buffers EarthlyBuffers) *bytes.Buffer {
	canvas := image.NewNRGBA(image.Rect(0, 0, config.Size, config.Size))
	earth, err := jpeg.Decode(bytes.NewReader(buffers.Earth))

	earthTexWidth := earth.Bounds().Dx()
	earthTexHeight := earth.Bounds().Dy()

	if err != nil {
		return new(bytes.Buffer)
	}

	var halfSize float64 = float64(config.Size) / 2.0
	var onePxSize = 1.0 / halfSize
	var twoPxSize = onePxSize * 2.0

	// A circle of 1.0 radius would touch the very edges of the image
	// In some cases we want to have border effects on the image so we make it a little smaller
	var circleRadius = 1.0
	if config.Radius > 0.0 {
		circleRadius = config.Radius
	}

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
		(latitudeCos*rollCos - latitudeSin*longitudeSin*rollSin), (latitudeSin*longitudeSin*rollCos + latitudeCos*rollSin), -(latitudeSin * longitudeCos),
		(longitudeCos * -rollSin), (longitudeCos * rollCos), (longitudeSin),
		(latitudeSin*rollCos + latitudeCos*longitudeSin*rollSin), (-latitudeCos*longitudeSin*rollCos + latitudeSin*rollSin), (latitudeCos * longitudeCos),
	}

	for py := 0; py < config.Size; py++ {
		for px := 0; px < config.Size; px++ {
			var r, g, b, a uint8 = 0, 0, 0, 0

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

			sphereX2 := matrix[0]*sphereX1 + matrix[1]*sphereY1 + matrix[2]*sphereZ1
			sphereY2 := matrix[3]*sphereX1 + matrix[4]*sphereY1 + matrix[5]*sphereZ1
			sphereZ2 := matrix[6]*sphereX1 + matrix[7]*sphereY1 + matrix[8]*sphereZ1

			projectedLongitude := math.Asin(sphereY2)
			projectedLatitude := math.Atan2(-sphereX2, sphereZ2)

			// Sample from the earth texture
			sampleX := earthTexWidth - 1 - (int(
				(projectedLatitude+math.Pi)*(float64(earthTexWidth)/(math.Pi*2)),
			)+earthTexWidth)%earthTexWidth

			sampleY := earthTexHeight - 1 - int(
				(projectedLongitude+math.Pi*0.5)*(float64(earthTexHeight)/math.Pi),
			)

			if sampleY < 0 {
				sampleY = 0
			}
			if sampleY >= earthTexHeight {
				sampleY = earthTexHeight
			}

			earthR, earthG, earthB, _ := earth.At(sampleX, sampleY).RGBA()

			r = uint8(earthR >> 8)
			g = uint8(earthG >> 8)
			b = uint8(earthB >> 8)

			// Mask it into a soft circle shape
			circleMask := 1.0 - math.Max(0.0, math.Min(1.0, (circleMagnitude-(1.0-twoPxSize))/twoPxSize))

			a = uint8(circleMask * 255.0)

			if config.Halo {
				haloMix := 255 - uint8(sphereZ1*255.0)
				haloR, haloG, haloB, _ := AlphaComposite(160, 160, 255, haloMix, 0, 0, 255, 255)
				r, g, b, _ = AlphaComposite(haloR, haloG, haloB, haloMix, r, g, b, 255)
			}

			r, g, b, a = AlphaComposite(r, g, b, a, config.Background[0], config.Background[1], config.Background[2], config.Background[3])

			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+0] = r
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+1] = g
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+2] = b
			canvas.Pix[(py-canvas.Rect.Min.Y)*canvas.Stride+(px-canvas.Rect.Min.X)*4+3] = a
		}
	}

	buffer := new(bytes.Buffer)
	png.Encode(buffer, canvas)

	return buffer
}
