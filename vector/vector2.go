package vector

type Vec2 struct {
	X float64
	Y float64
}

func (v1 Vec2) Dot(v2 Vec2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}
