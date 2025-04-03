package agents

import (
	"fmt"
	"math"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type PreoccupiedAgent struct {
	*ExtendedAgent
}

func CreatePreoccupiedAgent(server agent.IExposedServerFunctions[infra.IExtendedAgent], agentConfig AgentConfig, grid *infra.Grid) *PreoccupiedAgent {
	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)

	// Set Preoccupied-style attachment: high anxiety, low avoidance
	extendedAgent.Attachment = infra.Attachment{
		Anxiety:   randInRange(0.5, 1.0),
		Avoidance: randInRange(0.0, 0.5),
	}

	return &PreoccupiedAgent{
		ExtendedAgent: extendedAgent,
	}
}

func (pa *PreoccupiedAgent) AgentInitialised() {
	atch := pa.GetAttachment()
	fmt.Printf("Preoccupied Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", pa.GetID(), pa.GetAge(), atch.Anxiety, atch.Avoidance)
}

// preoccupied agent movement policy
// moves towards social network
func (pa *PreoccupiedAgent) GetTargetPosition(grid *infra.Grid) (infra.PositionVector, bool) {
	occupiedAgents := grid.GetAllOccupiedAgentPositions()
	//fmt.Printf("PreoccupiedAgent %v network: %v\n", pa.GetID(), pa.Network)

	var closestFriend infra.IExtendedAgent = nil
	minDist := math.MaxFloat32

	// Find closest friend
	for _, otherAgent := range occupiedAgents {
		if otherAgent.GetID() == pa.GetID() {
			continue // Skip self
		}
		if _, known := pa.Network[otherAgent.GetID()]; known {
			// friend so:
			dist := pa.Position.Dist(otherAgent.GetPosition())
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

// func (pa *PreoccupiedAgent) Move(grid *infra.Grid) {
// 	occupiedAgents := grid.GetAllOccupiedAgentPositions()
// 	//fmt.Printf("PreoccupiedAgent %v network: %v\n", pa.GetID(), pa.Network)

// 	var closestFriendID uuid.UUID
// 	var found bool
// 	minDist := math.MaxFloat32

// 	// Find closest friend
// 	for _, otherAgent := range occupiedAgents {
// 		if otherAgent.GetID() == pa.GetID() {
// 			continue // Skip self
// 		}
// 		if _, known := pa.Network[otherAgent.GetID()]; known {
// 			// friend so:
// 			dist := distance(pa.Position, otherAgent.GetPosition())
// 			if dist < minDist {
// 				minDist = dist
// 				closestFriendID = otherAgent.GetID()
// 				found = true
// 			}
// 		}
// 	}
// 	//fmt.Printf("PreoccupiedAgent %v found closest friend %v at distance %.2f\n", pa.GetID(), closestFriendID, minDist)

// 	if found {
// 		friend, ok := pa.Server.GetAgentMap()[closestFriendID]
// 		if ok {
// 			targetPos := friend.GetPosition()
// 			moveX := pa.Position[0] + getStep(pa.Position[0], targetPos[0])
// 			moveY := pa.Position[1] + getStep(pa.Position[1], targetPos[1])

// 			if moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY) {
// 				grid.UpdateAgentPosition(pa, moveX, moveY)
// 				pa.Position = [2]int{moveX, moveY}
// 				fmt.Printf("PreoccupiedAgent %v moved toward friend %v to (%d, %d)\n", pa.GetID(), closestFriendID, moveX, moveY)
// 				return
// 			}
// 		}
// 	}

// 	// Fallback: move randomly if no friends found
// 	newX, newY := grid.GetValidMove(pa.Position[0], pa.Position[1])
// 	grid.UpdateAgentPosition(pa, newX, newY)
// 	pa.Position = [2]int{newX, newY}
// 	fmt.Printf("PreoccupiedAgent %v fallback random move to (%d, %d)\n", pa.GetID(), newX, newY)
// }

//preoccupied agent pts protocol
//high probability of checking on other agents
//high probability of responding to other agents
