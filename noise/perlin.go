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

const GridSize uint = 20

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

func (noise *PerlinNoise) Get(x float64, y float64) (float64, error) {
    if x < 0.0 || x > 1.0 || y < 0.0 || y > 1.0 {
        return 0.0, fmt.Errorf(
            "Expected a x/y value between zero and one, got: %g, %g.", x, y,
        )    
    }

    scaledX := x * float64(noise.gridSize-1)
    scaledY := y * float64(noise.gridSize-1)

    x0, y0 := int(math.Floor(scaledX)), int(math.Floor(scaledY))
    x1, y1 := x0+1, y0+1

    x0 = x0 % int(noise.gridSize)
    x1 = x1 % int(noise.gridSize)
    y0 = y0 % int(noise.gridSize)
    y1 = y1 % int(noise.gridSize)

    g00 := noise.grid[y0*int(noise.gridSize)+x0]
    g01 := noise.grid[y1*int(noise.gridSize)+x0]
    g10 := noise.grid[y0*int(noise.gridSize)+x1]
    g11 := noise.grid[y1*int(noise.gridSize)+x1]

    xf, yf := scaledX-float64(x0), scaledY-float64(y0)
    
    d00 := vector.Vec2{X: xf, Y: yf}
    d01 := vector.Vec2{X: xf, Y: yf-1}
    d10 := vector.Vec2{X: xf-1, Y: yf}
    d11 := vector.Vec2{X: xf-1, Y: yf-1}

    s := g00.Dot(d00)
    t := g10.Dot(d10)
    u := g01.Dot(d01)
    v := g11.Dot(d11)

    sx := xf * xf * (3 - 2*xf)
    sy := yf * yf * (3 - 2*yf)

    a := s + sx*(t-s)
    b := u + sx*(v-u)
    value := a + sy*(b-a)

    return value, nil
}
