package agents

import (
	"fmt"
	"math"

	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type PreoccupiedAgent struct {
	*ExtendedAgent
}

func CreatePreoccupiedAgent(server infra.IServer, agentConfig AgentConfig, grid *infra.Grid) *PreoccupiedAgent {
	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)

	// Set Preoccupied-style attachment: high anxiety, low avoidance
	extendedAgent.Attachment = infra.Attachment{
		Anxiety:   randInRange(0.5, 1.0),
		Avoidance: randInRange(0.0, 0.5),
	}

	extendedAgent.PTW = infra.PTSParams{
		CheckProb: randInRange(0.5, 1.0), 
		ReplyProb: randInRange(0.5, 1.0),
		Alpha:     0.5, 
		Beta:      0.1, 
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
