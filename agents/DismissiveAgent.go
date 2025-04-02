package agents

import (
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type DismissiveAgent struct {
	*ExtendedAgent

}

func CreateDismissiveAgent(server agent.IExposedServerFunctions[infra.IExtendedAgent], agentConfig AgentConfig, grid *infra.Grid) *DismissiveAgent {
	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)
	
	// Set Dismissive-style attachment: low anxiety, high avoidance
	extendedAgent.Attachment = []float32{
		randInRange(0.0, 0.5),
		randInRange(0.5, 1.0),
	}

	return &DismissiveAgent{
		ExtendedAgent: extendedAgent,
	}
}

func (da *DismissiveAgent) AgentInitialised() {
	fmt.Printf("Dismissive Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", da.GetID(), da.GetAge(), da.GetAttachment()[0], da.GetAttachment()[1])
}

// dismissive agent movement policy
// moves away from social network
func (da *DismissiveAgent) Move(grid *infra.Grid) {
	occupiedAgents := grid.GetAllOccupiedAgentPositions()
	//fmt.Printf("DismissiveAgent %v network: %v\n", pa.GetID(), pa.Network)

	var closestFriendID uuid.UUID
	var found bool
	minDist := math.MaxFloat32

	// Find closest friend
	for _, otherAgent := range occupiedAgents {
		if otherAgent.GetID() == da.GetID() {
			continue // Skip self
		}
		if _, known := da.Network[otherAgent.GetID()]; known {
			// friend so:
			dist := distance(da.Position, otherAgent.GetPosition())
			if dist < minDist {
				minDist = dist
				closestFriendID = otherAgent.GetID()
				found = true
			}
		}
	}

	if found {
		friend, ok := da.Server.GetAgentMap()[closestFriendID]
		if ok {
			targetPos := friend.GetPosition()
			moveX := da.Position[0] - getStep(da.Position[0], targetPos[0])
			moveY := da.Position[1] - getStep(da.Position[1], targetPos[1])

			if moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY) {
				grid.UpdateAgentPosition(da, moveX, moveY)
				da.Position = [2]int{moveX, moveY}
				fmt.Printf("Dismissive Agent %v moved away from friend %v to (%d, %d)\n", da.GetID(), closestFriendID, moveX, moveY)
				return
			}
		}
	}
	

	// Fallback: move randomly if no friends found
	newX, newY := grid.GetValidMove(da.Position[0], da.Position[1])
	grid.UpdateAgentPosition(da, newX, newY)
	da.Position = [2]int{newX, newY}
	fmt.Printf("Dismissive Agent %v fallback random move to (%d, %d)\n", da.GetID(), newX, newY)
}



//dismissive agent pts protocol
// low probability of checking on other agents
// low probability of responding to other agents