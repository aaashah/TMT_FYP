package agents

import (
	"fmt"
	// "github.com/google/uuid"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type DismissiveAgent struct {
	*ExtendedAgent

}

func CreateDismissiveAgent(server agent.IExposedServerFunctions[infra.IExtendedAgent], agentConfig AgentConfig, grid *infra.Grid) *DismissiveAgent {

	extendedAgent := CreateExtendedAgent(server, agentConfig, grid)
	//extendedAgent := GetBaseAgents(server, agentConfig)

	return &DismissiveAgent{
		ExtendedAgent: extendedAgent,
	}
}

// dismissive agent movement policy
// moves away from social network
func (da *DismissiveAgent) Move(grid *infra.Grid) {
	if len(da.GetNetwork()) == 0 {
		// No connections – random move
		newX, newY := grid.GetValidMove(da.Position[0], da.Position[1])
		grid.UpdateAgentPosition(da, newX, newY)
		da.Position = [2]int{newX, newY}
		fmt.Printf("Dismissive Agent %v moved randomly to (%d, %d)\n", da.GetID(), newX, newY)
		return
	}

	closestFriend := da.FindClosestFriend()
	if closestFriend == nil {
		// Fallback random move
		newX, newY := grid.GetValidMove(da.Position[0], da.Position[1])
		grid.UpdateAgentPosition(da, newX, newY)
		da.Position = [2]int{newX, newY}
		fmt.Printf("Dismissive Agent %v moved randomly to (%d, %d)\n", da.GetID(), newX, newY)
		return
	}

	// Move one step away from the closest friend
	friendPos := closestFriend.Position
	moveX := da.Position[0] - getStep(friendPos[0], da.Position[0])
	moveY := da.Position[1] - getStep(friendPos[1], da.Position[1])

	if moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY) {
		grid.UpdateAgentPosition(da, moveX, moveY)
		da.Position = [2]int{moveX, moveY}
		fmt.Printf("Dismissive Agent %v moved away from friend %v to (%d, %d)\n", da.GetID(), closestFriend.GetID(), moveX, moveY)
	} else {
		newX, newY := grid.GetValidMove(da.Position[0], da.Position[1])
		grid.UpdateAgentPosition(da, newX, newY)
		da.Position = [2]int{newX, newY}
		fmt.Printf("Dismissive Agent %v tried to avoid friend %v but was blocked, moved randomly to (%d, %d)\n", da.GetID(), closestFriend.GetID(), newX, newY)
	}
}



//dismissive agent pts protocol
// low probability of checking on other agents
// low probability of responding to other agents