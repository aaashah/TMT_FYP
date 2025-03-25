package agents

import (

	// "github.com/google/uuid"
	"fmt"
	//"math"

	//"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	//gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
)

type PreoccupiedAgent struct {
	*ExtendedAgent

}

func CreatePreoccupiedAgent(server infra.IServer , agentConfig AgentConfig, grid *infra.Grid) *PreoccupiedAgent {
	
	extendedAgent := CreateExtendedAgents(server, agentConfig, grid)

	return &PreoccupiedAgent{
		ExtendedAgent: extendedAgent,
	}
}

// preoccupied agent movement policy
// moves towards social network
func (pa *PreoccupiedAgent) Move(grid *infra.Grid) {
	if len(pa.Network) == 0 {
		// No social connections, move randomly
		newX, newY := grid.GetValidMove(pa.Position[0], pa.Position[1])
		grid.UpdateAgentPosition(pa, newX, newY)
		pa.Position = [2]int{newX, newY}
		fmt.Printf("Agent %v moved randomly to (%d, %d)\n", pa.GetID(), newX, newY)
		return
	}

// 	Find closest friend in the social network
	closestFriend := pa.FindClosestFriend()
	if closestFriend == nil {
		// No valid friends found, move randomly
		newX, newY := grid.GetValidMove(pa.Position[0], pa.Position[1])
		grid.UpdateAgentPosition(pa, newX, newY)
		pa.Position = [2]int{newX, newY}
		fmt.Printf("Agent %v moved randomly to (%d, %d)\n", pa.GetID(), newX, newY)
		return
	}

	// Move one step towards closest friend
	targetPos := closestFriend.Position
	moveX := pa.Position[0] + getStep(pa.Position[0], targetPos[0])
	moveY := pa.Position[1] + getStep(pa.Position[1], targetPos[1])

	// Ensure movement stays within grid bounds and doesn't move to occupied space
	if moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY){
		grid.UpdateAgentPosition(pa, moveX, moveY)
		pa.Position = [2]int{moveX, moveY}
		fmt.Printf("Agent %v moved towards friend %v to (%d, %d)\n", pa.GetID(), closestFriend.GetID(), moveX, moveY)
	} else {
		// Blocked — fallback to random valid move
		newX, newY := grid.GetValidMove(pa.Position[0], pa.Position[1])
		grid.UpdateAgentPosition(pa, newX, newY)
		pa.Position = [2]int{newX, newY}
		fmt.Printf("Agent %v tried to move toward friend %v but was blocked, moved randomly to (%d, %d)\n", pa.GetID(), closestFriend.GetID(), newX, newY)
	}
}


//preoccupied agent pts protocol
//high probability of checking on other agents
//high probability of responding to other agents