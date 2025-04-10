package agents

import (
	"fmt"
	// "github.com/google/uuid"

	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type SecureAgent struct {
	*ExtendedAgent
}

func CreateSecureAgent(server infra.IServer, agentConfig AgentConfig, grid *infra.Grid) *SecureAgent {
	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)

	// Set Secure-style attachment: low anxiety, low avoidance
	extendedAgent.Attachment = infra.Attachment{
		Anxiety:   randInRange(0.0, 0.5),
		Avoidance: randInRange(0.0, 0.5),
	}

	extendedAgent.PTW = infra.PTSParams{
		CheckProb: randInRange(0.0, 0.5), 
		ReplyProb: randInRange(0.5, 1.0),
		Alpha:     0.5, 
		Beta:      0.1, 
	}

	return &SecureAgent{
		ExtendedAgent: extendedAgent,
	}
}

func (sa *SecureAgent) AgentInitialised() {
	atch := sa.GetAttachment()
	fmt.Printf("Secure Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", sa.GetID(), sa.GetAge(), atch.Anxiety, atch.Avoidance)
}

// Secure agent movement policy
// moves randomly
func (sa *SecureAgent) GetTargetPosition(grid *infra.Grid) (infra.PositionVector, bool) {
	return infra.PositionVector{}, false
}


