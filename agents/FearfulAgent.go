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

	// Set Fearful-style attachment: high anxiety, high avoidance
	extendedAgent.Attachment = []float32{
		randInRange(0.5, 1.0),
		randInRange(0.5, 1.0),
	}

	return &FearfulAgent{
		ExtendedAgent: extendedAgent,
	}
}
func (fa *FearfulAgent) AgentInitialised() {
	fmt.Printf("Fearful Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", fa.GetID(), fa.GetAge(), fa.GetAttachment()[0], fa.GetAttachment()[1])
}

// Fearful agent movement policy
// moves towards those not in social network
func (fa *FearfulAgent) Move(grid *infra.Grid) {
	occupied := grid.GetAllOccupiedAgentPositions()

	var closestStrangerID uuid.UUID
	var found bool
	minDist := math.MaxFloat32

	for _, otherAgent := range occupied {
		if otherAgent.GetID() == fa.GetID() {
			continue // Skip self
		}
		if _, known := fa.Network[otherAgent.GetID()]; known {
			continue // Skip friends
		}
	
		dist := distance(fa.Position, otherAgent.GetPosition())
		if dist < minDist {
			minDist = dist
			closestStrangerID = otherAgent.GetID()
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