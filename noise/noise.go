package noise

type Noise interface {
	Get(x float64, y float64) (float64, error)
}
