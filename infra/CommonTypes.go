package infra

import "math"

type PositionVector struct {
	X int
	Y int
}

func (p1 PositionVector) Dist(p2 PositionVector) float64 {
	deltaX := float64(p1.X - p2.X)
	deltaY := float64(p1.Y - p2.Y)
	return math.Sqrt((deltaX * deltaX) + (deltaY * deltaY))
}
