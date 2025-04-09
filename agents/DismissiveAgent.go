package agents

import (
	"fmt"
	"math"

	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type DismissiveAgent struct {
	*ExtendedAgent
}

func CreateDismissiveAgent(server infra.IServer, agentConfig AgentConfig, grid *infra.Grid) *DismissiveAgent {
	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)

	// Set Dismissive-style attachment: low anxiety, high avoidance
	extendedAgent.Attachment = infra.Attachment{
		Anxiety:   randInRange(0.0, 0.5),
		Avoidance: randInRange(0.5, 1.0),
	}

	return &DismissiveAgent{
		ExtendedAgent: extendedAgent,
	}
}

func (da *DismissiveAgent) AgentInitialised() {
	atch := da.GetAttachment()
	fmt.Printf("Dismissive Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", da.GetID(), da.GetAge(), atch.Anxiety, atch.Avoidance)
}

// dismissive agent movement policy
// moves away from social network
func (da *DismissiveAgent) GetTargetPosition(grid *infra.Grid) (infra.PositionVector, bool) {
	occupiedAgents := grid.GetAllOccupiedAgentPositions()
	//fmt.Printf("DismissiveAgent %v network: %v\n", pa.GetID(), pa.Network)

	var closestFriend infra.IExtendedAgent = nil
	minDist := math.Inf(1)

	// Find closest friend
	for _, otherAgent := range occupiedAgents {
		if otherAgent.GetID() == da.GetID() {
			continue // Skip self
		}
		if _, known := da.Network[otherAgent.GetID()]; known {
			// friend so:
			dist := da.Position.Dist(otherAgent.GetPosition())
			if dist < minDist {
				minDist = dist
				closestFriend = otherAgent
			}
		}
	}

	if closestFriend == nil {
		return infra.PositionVector{}, false
	}

	return closestFriend.GetPosition(), true
}

// func (da *DismissiveAgent) Move(grid *infra.Grid) {
// 	occupiedAgents := grid.GetAllOccupiedAgentPositions()
// 	//fmt.Printf("DismissiveAgent %v network: %v\n", pa.GetID(), pa.Network)

// 	var closestFriendID uuid.UUID
// 	var found bool
// 	minDist := math.MaxFloat32

// 	// Find closest friend
// 	for _, otherAgent := range occupiedAgents {
// 		if otherAgent.GetID() == da.GetID() {
// 			continue // Skip self
// 		}
// 		if _, known := da.Network[otherAgent.GetID()]; known {
// 			// friend so:
// 			dist := da.Position.Dist(otherAgent.GetPosition())
// 			if dist < minDist {
// 				minDist = dist
// 				closestFriendID = otherAgent.GetID()
// 				found = true
// 			}
// 		}
// 	}

// 	if found {
// 		friend, ok := da.Server.GetAgentMap()[closestFriendID]
// 		if ok {
// 			targetPos := friend.GetPosition()
// 			moveX := da.Position.X - getStep(da.Position.X, targetPos.X)
// 			moveY := da.Position.Y - getStep(da.Position.Y, targetPos.Y)

// 			if moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY) {
// 				grid.UpdateAgentPosition(da, moveX, moveY)
// 				da.Position = infra.PositionVector{X: moveX, Y: moveY}
// 				fmt.Printf("Dismissive Agent %v moved away from friend %v to (%d, %d)\n", da.GetID(), closestFriendID, moveX, moveY)
// 				return
// 			}
// 		}
// 	}

// 	// Fallback: move randomly if no friends found
// 	newX, newY := grid.GetValidMove(da.Position.X, da.Position.Y)
// 	grid.UpdateAgentPosition(da, newX, newY)
// 	da.Position = infra.PositionVector{X: newX, Y: newY}
// 	fmt.Printf("Dismissive Agent %v fallback random move to (%d, %d)\n", da.GetID(), newX, newY)
// }

//dismissive agent pts protocol
// low probability of checking on other agents
// low probability of responding to other agents
