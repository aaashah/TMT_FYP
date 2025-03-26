package agents

import (
	"fmt"
	"math"

	// "github.com/google/uuid"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/google/uuid"

	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type FearfulAgent struct {
	*ExtendedAgent

}

func CreateFearfulAgent(server agent.IExposedServerFunctions[infra.IExtendedAgent] , agentConfig AgentConfig, grid *infra.Grid) *FearfulAgent {
	
	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)

	return &FearfulAgent{
		ExtendedAgent: extendedAgent,
	}
}

// Fearful agent movement policy
// moves towards those not in social network
func (fa *FearfulAgent) Move(grid *infra.Grid) {
	occupied := grid.GetAllOccupiedAgentPositions()

	var closestStrangerID uuid.UUID
	var found bool
	minDist := math.MaxFloat32

	for _, other := range occupied {
		if other.GetID() == fa.GetID() {
			continue
		}
		if _, known := fa.Network[other.GetID()]; known {
			continue
		}
	
		dist := distance(fa.Position, other.GetPosition())
		if dist < minDist {
			minDist = dist
			closestStrangerID = other.GetID()
			found = true
		}
	}

	if found {
		stranger, ok := fa.Server.GetAgentMap()[closestStrangerID]
		if ok {
			targetPos := stranger.GetPosition()
			moveX := fa.Position[0] + getStep(fa.Position[0], targetPos[0])
			moveY := fa.Position[1] + getStep(fa.Position[1], targetPos[1])
	
			if moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY) {
				grid.UpdateAgentPosition(fa, moveX, moveY)
				fa.Position = [2]int{moveX, moveY}
				fmt.Printf("FearfulAgent %v moved toward stranger %v to (%d, %d)\n", fa.GetID(), closestStrangerID, moveX, moveY)
				return
			}
		}
	}

	// Fallback: move randomly if no strangers found
	newX, newY := grid.GetValidMove(fa.Position[0], fa.Position[1])
	grid.UpdateAgentPosition(fa, newX, newY)
	fa.Position = [2]int{newX, newY}
	fmt.Printf("FearfulAgent %v fallback random move to (%d, %d)\n", fa.GetID(), newX, newY)
}



//fearful agent pts protocol
//low probability of checking
// high probability of responding