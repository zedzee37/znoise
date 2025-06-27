package noise

import (
	"math"
	"math/rand"
	"znoise/vector"
)

type PerlinNoise struct {
	Octaves int
	Lacunarity float64
	Frequency float64
	Persistence float64
	Gain float64
	grids [][]vector.Vec2
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

func NewPerlinNoise(
	seed int64,
	octaves int,
	lacunarity float64,
	frequency float64,
	persistence float64,
	gain float64,
) PerlinNoise {
	randSource := rand.NewSource(seed)
	rng := rand.New(randSource)

	grids := make([][]vector.Vec2, octaves)

	for i := range octaves {
		grid := make([]vector.Vec2, GridSize*GridSize)
		for j := range GridSize*GridSize {
			grid[j] = generateRandomVector(rng)		
		}
		grids[i] = grid
	}

	return PerlinNoise {
		octaves,
		lacunarity,
		frequency,
		persistence,
		gain,
		grids,
		GridSize,	
		rng,
	}
} 

func lerp(a float64, b float64, p float64) float64 {
	return (1 - p)*a + p*b
}

func normalize(n float64) float64 {
	return (1.0 + n) / 2.0
}

func fade(t float64) float64 {
    return t * t * t * (t * (t * 6 - 15) + 10)
}

func (noise *PerlinNoise) Get(x float64, y float64) (float64, error) {
	value := 0.0	

	for i := range noise.Octaves {
		frequency := noise.Frequency * math.Pow(noise.Lacunarity, float64(i))	
		
		adjustedX := x * frequency
		adjustedY := y * frequency

		noiseValue, err := getPerlin(noise.grids[i], adjustedX, adjustedY, noise.gridSize)

		if err != nil {
			return 0.0, err
		}

		amplitude := math.Pow(noise.Persistence, float64(i))

		noiseValue *= amplitude
		value += noiseValue
	}

	value = normalize(value)
	if value > 1.0 {
		value = 1.0
	} else if value < 0.0 {
		value = 0.0
	}
	return value, nil
}

func getPerlin(grid []vector.Vec2, x float64, y float64, gridSize uint) (float64, error) {
    scaledX := x * float64(gridSize-1)
    scaledY := y * float64(gridSize-1)

    x0, y0 := int(math.Floor(scaledX)), int(math.Floor(scaledY))
    x1, y1 := x0+1, y0+1

    fracX := scaledX - float64(x0)
    fracY := scaledY - float64(y0)

    x0 = x0 % int(gridSize)
    x1 = x1 % int(gridSize)
    y0 = y0 % int(gridSize)
    y1 = y1 % int(gridSize)

    v00 := grid[y0*int(gridSize) + x0]    
    v10 := grid[y0*int(gridSize) + x1]    
    v01 := grid[y1*int(gridSize) + x0]    
    v11 := grid[y1*int(gridSize) + x1]    

    offset00 := vector.Vec2{X: fracX, Y: fracY}
    offset10 := vector.Vec2{X: fracX - 1.0, Y: fracY}
    offset01 := vector.Vec2{X: fracX, Y: fracY - 1.0}
    offset11 := vector.Vec2{X: fracX - 1.0, Y: fracY - 1.0}

    u := fade(fracX)
    v := fade(fracY)

    dot00 := offset00.Dot(v00)
    dot10 := offset10.Dot(v10)
    dot01 := offset01.Dot(v01)
    dot11 := offset11.Dot(v11)

    top := lerp(dot00, dot10, u)
    bottom := lerp(dot01, dot11, u)
    value := lerp(top, bottom, v)

    return value, nil
}
