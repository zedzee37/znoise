package main

import (
	"image/png"
	"os"
	"znoise/image"
	"znoise/noise"
)

func main() {
	perlinNoise := noise.NewPerlinNoise(100, 3, 2.0, 1.0, 0.5, 1.0)
	img, err := image.CreateNoiseImage(&perlinNoise, 2000, 2000)

	if err != nil {
		panic(err)
	}

	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}
