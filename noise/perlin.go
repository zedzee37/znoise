package noise

import (
	"fmt"
	"math"
	"math/rand"
	"znoise/vector"
)

type PerlinNoise struct {
	Seed int64
	Octaves int
	Lacunarity float32
	Gain float32
	grid []vector.Vec2
	gridSize uint
	rng *rand.Rand
}

const GridSize uint = 256

func generateRandomVector(rng *rand.Rand) vector.Vec2 {
	angle := rng.Float64() * (2 * math.Pi)
	return vector.Vec2{
		X: math.Cos(angle),
		Y: math.Sin(angle),
	}
}

func NewPerlinNoise(seed int64, octaves int, lacunarity float32, gain float32) PerlinNoise {
	randSource := rand.NewSource(seed)
	rng := rand.New(randSource)

	grid := make([]vector.Vec2, GridSize*GridSize)

	for i := range GridSize*GridSize {
		grid[i] = generateRandomVector(rng)		
	}

	return PerlinNoise {
		seed,
		octaves,
		lacunarity,
		gain,
		grid,
		GridSize,	
		rng,
	}
} 

func (noise PerlinNoise) Get(x float64, y float64) (float64, error) {
	if x < 0.0 || x > 1.0 || y < 0.0 || y > 1.0 {
		return 0.0, fmt.Errorf(
			"Expected a x/y value between zero and one, got: %g, %g.", x, y,
		)	
	}

	scaledX := x * float64(noise.gridSize)
	scaledY := y * float64(noise.gridSize)

	leftX, topY := int(math.Floor(scaledX)), int(math.Floor(scaledY))
	rightX := leftX + 1
	bottomY := topY + 1	

	v1 := noise.grid[topY * int(noise.gridSize) + leftX]
	v2 := noise.grid[topY * int(noise.gridSize) + rightX]
	v3 := noise.grid[bottomY * int(noise.gridSize) + leftX]
	v4 := noise.grid[bottomY * int(noise.gridSize) + rightX]
	
	coords := vector.Vec2{
		X: x,
		Y: y,
	}

	v1Offset := vector.Vec2{
		X: float64(leftX) / 512,
		Y: float64(topY) / 512,
	}.Sub(coords)
	v2Offset := vector.Vec2{
		X: float64(rightX) / 512,
		Y: float64(topY) / 512,
	}.Sub(coords)
	v3Offset := vector.Vec2{
		X: float64(leftX) / 512,
		Y: float64(bottomY) / 512,
	}.Sub(coords)
	v4Offset := vector.Vec2{
		X: float64(rightX) / 512,
		Y: float64(bottomY) / 512,
	}.Sub(coords)

	v1Dot := v1.Dot(v1Offset)
	v2Dot := v2.Dot(v2Offset)
	v3Dot := v3.Dot(v3Offset)
	v4Dot := v4.Dot(v4Offset)

	return 1.0, nil
}
