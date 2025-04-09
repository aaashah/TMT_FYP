package infra

import (
	//"math"
	"math/rand"
	"sync"
)

type Grid struct {
	Width      int
	Height     int
	positions  map[PositionVector]IExtendedAgent
	Tombstones []PositionVector
	Temples    []PositionVector
	mutex      sync.Mutex
}

func NewGrid(width, height int) *Grid {
	return &Grid{
		Width:      width,
		Height:     height,
		positions:  make(map[PositionVector]IExtendedAgent),
		Tombstones: []PositionVector{},
		Temples:    []PositionVector{},
	}
}

// Check if a cell is occupied
func (g *Grid) IsOccupied(x, y int) bool {
	pos := PositionVector{X: x, Y: y}
	if _, exists := g.positions[pos]; exists {
		return true
	}

	target := PositionVector{X: x, Y: y}
	for _, t := range g.Tombstones {
		if t == target {
			return true
		}
	}
	for _, temple := range g.Temples {
		if temple == target {
			return true
		}
	}
	return false
}

// Place a tombstone at an agent's last known position
func (g *Grid) PlaceTombstone(x, y int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Tombstones = append(g.Tombstones, PositionVector{X: x, Y: y}) // Mark the position as a tombstone
}

func (g *Grid) PlaceTemple(x, y int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Temples = append(g.Temples, PositionVector{X: x, Y: y}) // Mark the position as a temple
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
	return x, y // Stay in place if no move is available
}

// Update agent position on the grid
func (g *Grid) UpdateAgentPosition(agent IExtendedAgent, newPos PositionVector) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Ensure the new position is not already occupied OR a tombstone
	if g.IsOccupied(newPos.X, newPos.Y) {
		//fmt.Printf(" Agent %v tried to move onto an occupied cell (%d, %d). Movement canceled.\n", agent.GetID(), newX, newY)
		return
	}

	// Remove from old position
	oldPos := agent.GetPosition()
	delete(g.positions, oldPos)

	// Update new position
	g.positions[newPos] = agent
}

func (g *Grid) GetAllOccupiedAgentPositions() map[PositionVector]IExtendedAgent {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Return a copy to prevent accidental modification
	copyMap := make(map[PositionVector]IExtendedAgent)
	for pos, agent := range g.positions {
		copyMap[pos] = agent
	}
	return copyMap
}
