package agents

import (
	"fmt"
	"math"

	// "github.com/google/uuid"

	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type FearfulAgent struct {
	*ExtendedAgent
}

func CreateFearfulAgent(server infra.IServer, agentConfig AgentConfig, grid *infra.Grid) *FearfulAgent {
	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)

	// Set Fearful-style attachment: high anxiety, high avoidance
	extendedAgent.Attachment = infra.Attachment{
		Anxiety:   randInRange(0.5, 1.0),
		Avoidance: randInRange(0.5, 1.0),
	}

	return &FearfulAgent{
		ExtendedAgent: extendedAgent,
	}
}
func (fa *FearfulAgent) AgentInitialised() {
	atch := fa.GetAttachment()
	fmt.Printf("Fearful Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", fa.GetID(), fa.GetAge(), atch.Anxiety, atch.Avoidance)
}

// Fearful agent movement policy
// moves towards those not in social network
func (fa *FearfulAgent) GetTargetPosition(grid *infra.Grid) (infra.PositionVector, bool) {
	occupied := grid.GetAllOccupiedAgentPositions()

	var closestStranger infra.IExtendedAgent = nil
	minDist := math.MaxFloat32

	for _, otherAgent := range occupied {
		if otherAgent.GetID() == fa.GetID() {
			continue // Skip self
		}
		if _, known := fa.Network[otherAgent.GetID()]; known {
			continue // Skip friends
		}

		dist := fa.Position.Dist(otherAgent.GetPosition())
		if dist < minDist {
			minDist = dist
			closestStranger = otherAgent
		}
	}

	if closestStranger == nil {
		return infra.PositionVector{}, false
	}

	return closestStranger.GetPosition(), true

}

// func (fa *FearfulAgent) Move(grid *infra.Grid) {
// 	occupied := grid.GetAllOccupiedAgentPositions()

// 	var closestStrangerID uuid.UUID
// 	var found bool
// 	minDist := math.MaxFloat32

// 	for _, otherAgent := range occupied {
// 		if otherAgent.GetID() == fa.GetID() {
// 			continue // Skip self
// 		}
// 		if _, known := fa.Network[otherAgent.GetID()]; known {
// 			continue // Skip friends
// 		}

// 		dist := fa.Position.Dist(otherAgent.GetPosition())
// 		if dist < minDist {
// 			minDist = dist
// 			closestStrangerID = otherAgent.GetID()
// 			found = true
// 		}
// 	}

// 	if found {
// 		stranger, ok := fa.Server.GetAgentMap()[closestStrangerID]
// 		if ok {
// 			targetPos := stranger.GetPosition()
// 			moveX := fa.Position.X + getStep(fa.Position.X, targetPos.X)
// 			moveY := fa.Position.Y + getStep(fa.Position.Y, targetPos.Y)

// 			if moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY) {
// 				grid.UpdateAgentPosition(fa, moveX, moveY)
// 				fa.Position = infra.PositionVector{X: moveX, Y: moveY}
// 				fmt.Printf("FearfulAgent %v moved toward stranger %v to (%d, %d)\n", fa.GetID(), closestStrangerID, moveX, moveY)
// 				return
// 			}
// 		}
// 	}

// 	// Fallback: move randomly if no strangers found
// 	newX, newY := grid.GetValidMove(fa.Position.X, fa.Position.Y)
// 	grid.UpdateAgentPosition(fa, newX, newY)
// 	fa.Position = infra.PositionVector{X: newX, Y: newY}
// 	fmt.Printf("FearfulAgent %v fallback random move to (%d, %d)\n", fa.GetID(), newX, newY)
// }

//fearful agent pts protocol
//low probability of checking
// high probability of responding
