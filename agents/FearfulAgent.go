package agents

import (
	"fmt"
	"math"

	// "github.com/google/uuid"
	//"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type FearfulAgent struct {
	*ExtendedAgent

}

func CreateFearfulAgent(server infra.IServer , agentConfig AgentConfig, grid *infra.Grid) *FearfulAgent {
	
	extendedAgent := CreateExtendedAgents(server, agentConfig, grid)

	return &FearfulAgent{
		ExtendedAgent: extendedAgent,
	}
}

// Fearful agent movement policy
// moves towards those not in social network
func (fa *FearfulAgent) Move(grid *infra.Grid) {
	allAgents := fa.Server.GetAgentMap()
	var closestStranger *ExtendedAgent
	minDist := math.MaxFloat64

	for id, agentInterface := range allAgents {
		// Skip self and agents already in social network
		if id == fa.GetID() || fa.Network[id] > 0 {
			continue
		}

		stranger, ok := agentInterface.(*ExtendedAgent)
		if !ok {
			continue
		}

		dist := distance(fa.Position, stranger.Position)
		if dist < minDist {
			minDist = dist
			closestStranger = stranger
		}
	}

	if closestStranger == nil {
		// No strangers found, move randomly
		newX, newY := grid.GetValidMove(fa.Position[0], fa.Position[1])
		grid.UpdateAgentPosition(fa, newX, newY)
		fa.Position = [2]int{newX, newY}
		fmt.Printf("FearfulAgent %v moved randomly to (%d, %d)\n", fa.GetID(), newX, newY)
		return
	}

	// Move one step toward closest stranger
	targetPos := closestStranger.Position
	moveX := fa.Position[0] + getStep(fa.Position[0], targetPos[0])
	moveY := fa.Position[1] + getStep(fa.Position[1], targetPos[1])

	// If move is valid and unoccupied
	if moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY) {
		grid.UpdateAgentPosition(fa, moveX, moveY)
		fa.Position = [2]int{moveX, moveY}
		fmt.Printf("FearfulAgent %v moved toward stranger %v to (%d, %d)\n", fa.GetID(), closestStranger.GetID(), moveX, moveY)
	} else {
		// Fallback: move randomly
		newX, newY := grid.GetValidMove(fa.Position[0], fa.Position[1])
		grid.UpdateAgentPosition(fa, newX, newY)
		fa.Position = [2]int{newX, newY}
		fmt.Printf("FearfulAgent %v fallback random move to (%d, %d)\n", fa.GetID(), newX, newY)
	}
}



//fearful agent pts protocol
//low probability of checking
// high probability of responding