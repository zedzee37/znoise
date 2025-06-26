package image

import (
	"image"
	"image/color"
	"znoise/noise"
)

func CreateNoiseImage(noise noise.Noise, width uint, height uint) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rectangle{
		image.Point{
			X: 0,
			Y: 0,
		}, 
		image.Point{
			X: int(width),
			Y: int(height),
		},
	})
	
	for x := range width {
		for y := range height {
			scaledX := float64(x) / float64(width)
			scaledY := float64(y) / float64(height)
			noiseValue, err := noise.Get(scaledX, scaledY)			

			if err != nil {
				return nil, err
			}

			scaledNoise := uint8(noiseValue * 255)

			img.Set(int(x), int(y), color.RGBA{
				R: scaledNoise,
				G: scaledNoise,
				B: scaledNoise,
				A: uint8(255),
			})
		}
	}

	return img, nil
}
