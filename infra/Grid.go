package infra

import (
	//"math"
	"math/rand"
	"sync"
)

type Grid struct {
    Width  int
    Height int
	positions     map[[2]int]IExtendedAgent
	Tombstones map[[2]int]bool
	Temples map[[2]int]bool
    mutex     sync.Mutex
}

func NewGrid(width, height int) *Grid {
	return &Grid{
		Width: width,
		Height: height,
		positions: make(map[[2]int]IExtendedAgent),
		Tombstones: make(map[[2]int]bool),
		Temples: make(map[[2]int]bool),
	}
}

// Check if a cell is occupied
func (g *Grid) IsOccupied(x, y int) bool {
	_, exists := g.positions[[2]int{x, y}]
	_, isTombstone := g.Tombstones[[2]int{x, y}] // ✅ Correctly using both values
	_, isTemple := g.Temples[[2]int{x, y}] // ✅ Correctly using both values

	return exists || isTombstone || isTemple// ✅ Now `isTombstone` is properly used
}

// Place a tombstone at an agent's last known position
func (g *Grid) PlaceTombstone(x, y int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Tombstones[[2]int{x, y}] = true  // ✅ Mark the position as a tombstone
}

func (g *Grid) PlaceTemple(x, y int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Temples[[2]int{x, y}] = false  // ✅ Mark the position as a temple
}

// Get a valid move for an agent
func (g *Grid) GetValidMove(x, y int) (int, int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	moves := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })

	for _, move := range moves {
		newX, newY := x+move[0], y+move[1]
		if newX > 0 && newX <= g.Width && newY > 0 && newY <= g.Height && !g.IsOccupied(newX, newY) {
			return newX, newY
		}
	}
	return x, y  // Stay in place if no move is available
}


// Update agent position on the grid
func (g *Grid) UpdateAgentPosition(agent IExtendedAgent, newX, newY int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Ensure the new position is not already occupied OR a tombstone
	if g.IsOccupied(newX, newY) {
		//fmt.Printf("⚠️ Agent %v tried to move onto an occupied cell (%d, %d). Movement canceled.\n", agent.GetID(), newX, newY)
		return
	}

	// Remove from old position
	oldPos := agent.GetPosition()
	delete(g.positions, [2]int{oldPos[0], oldPos[1]})

	// Update new position
	g.positions[[2]int{newX, newY}] = agent
}