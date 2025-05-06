package agents

import (
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/MattSScott/TMT_SOMAS/infra"
)

type FearfulAgent struct {
	*ExtendedAgent
}

func CreateFearfulAgent(server infra.IServer, parent1ID uuid.UUID, parent2ID uuid.UUID) *FearfulAgent {
	worldview := infra.NewWorldview(byte(0b00))
	extendedAgent := CreateExtendedAgent(server, parent1ID, parent2ID, worldview)

	// Set Fearful-style attachment: high anxiety, high avoidance
	extendedAgent.attachment = infra.Attachment{
		Anxiety:   randInRange(0.5, 1.0),
		Avoidance: randInRange(0.5, 1.0),
		Type:      infra.FEARFUL,
	}
	// these ranges to be tweaked
	extendedAgent.PTW = infra.PTSParams{
		CheckProb: randInRange(0.5, 1.0),
		ReplyProb: randInRange(0.0, 0.5),
		Alpha:     randInRange(0.0, 0.5),
		Beta:      randInRange(0.5, 1.0),
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
	minDist := math.Inf(1)

	for _, otherAgent := range occupied {
		if otherAgent.GetID() == fa.GetID() {
			continue // Skip self
		}
		if _, known := fa.network[otherAgent.GetID()]; known {
			continue // Skip friends
		}

		dist := fa.position.Dist(otherAgent.GetPosition())
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
