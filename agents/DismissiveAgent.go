package agents

import (
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/aaashah/TMT_Attachment/infra"
)

type DismissiveAgent struct {
	*ExtendedAgent
}

func CreateDismissiveAgent(server infra.IServer, parent1ID uuid.UUID, parent2ID uuid.UUID, worldview uint32) *DismissiveAgent {
	extendedAgent := CreateExtendedAgent(server, parent1ID, parent2ID, worldview)

	// Set Dismissive-style attachment: low anxiety, high avoidance
	extendedAgent.Attachment = infra.Attachment{
		Anxiety:   randInRange(0.0, 0.5),
		Avoidance: randInRange(0.5, 1.0),
		Type:      infra.DISMISSIVE,
	}
	// these ranges to be tweaked
	extendedAgent.PTW = infra.PTSParams{
		CheckProb: randInRange(0.0, 0.5),
		ReplyProb: randInRange(0.0, 0.5),
		Alpha:     randInRange(0.0, 0.5),
		Beta:      randInRange(0.0, 0.5),
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
		if _, known := da.network[otherAgent.GetID()]; known {
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
