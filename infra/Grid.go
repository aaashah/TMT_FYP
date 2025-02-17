package infra

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