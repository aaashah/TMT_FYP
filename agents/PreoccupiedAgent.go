package agents

import (
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/aaashah/TMT_Attachment/infra"
)

type PreoccupiedAgent struct {
	*ExtendedAgent
}

func CreatePreoccupiedAgent(server infra.IServer, parent1ID uuid.UUID, parent2ID uuid.UUID, worldview uint32) *PreoccupiedAgent {
	extendedAgent := CreateExtendedAgent(server, parent1ID, parent2ID, worldview)

	// Set Preoccupied-style attachment: high anxiety, low avoidance
	extendedAgent.Attachment = infra.Attachment{
		Anxiety:   randInRange(0.5, 1.0),
		Avoidance: randInRange(0.0, 0.5),
		Type:      "Preoccupied",
	}
	// these ranges to be tweaked
	extendedAgent.PTW = infra.PTSParams{
		CheckProb: randInRange(0.5, 1.0),
		ReplyProb: randInRange(0.5, 1.0),
		Alpha:     randInRange(0.5, 1.0),
		Beta:      randInRange(0.5, 1.0),
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
	minDist := math.Inf(1)

	// Find closest friend
	for _, otherAgent := range occupiedAgents {
		if otherAgent.GetID() == pa.GetID() {
			continue // Skip self
		}
		if _, known := pa.network[otherAgent.GetID()]; known {
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
