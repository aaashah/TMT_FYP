package infra

import (
	"math"
)

type Grid struct {
    Width  int
    Height int
    //Cells  map[Position][]*agents.ExtendedAgent
    //mu     sync.Mutex
}

func CreateGrid(width, height int) *Grid {
	return &Grid{
		Width: width,
		Height: height,
	}
}

func Distance(a [2]int, b [2]int) float64 {
	return math.Sqrt(float64((a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1])))
}