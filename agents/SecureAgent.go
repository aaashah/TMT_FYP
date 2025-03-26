package agents

import (
	"fmt"
	// "github.com/google/uuid"
	// "math/rand"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type SecureAgent struct {
	*ExtendedAgent

}

func CreateSecureAgent(server agent.IExposedServerFunctions[infra.IExtendedAgent] , agentConfig AgentConfig, grid *infra.Grid) *SecureAgent {
	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)

	return &SecureAgent{
		ExtendedAgent: extendedAgent,
	}
}

// Secure agent movement policy
// moves randomly
func (sa *SecureAgent) Move(grid *infra.Grid) {
	newX, newY := grid.GetValidMove(sa.Position[0], sa.Position[1]) // Get a valid move
	grid.UpdateAgentPosition(sa, newX, newY)    // Update position in the grid
	sa.Position = [2]int{newX, newY}             // Assign new position
	fmt.Printf("Secure Agent %v moved to (%d, %d)\n", sa.GetID(), newX, newY)
}

//secure agent pts protocol
//low probability of checking on other agents
//high probability of responding to other agents