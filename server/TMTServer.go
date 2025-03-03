package server

import (
	"fmt"
	"log"
	"math/rand"

	//"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"

	agents "github.com/aaashah/TMT_Attachment/agents"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
) 

type TMTServer struct {
	*server.BaseServer[infra.IExtendedAgent]

	//agentInfoList []infra.IExtendedAgent
	//mu     sync.Mutex
	context string
	ActiveAgents map[uuid.UUID]*agents.ExtendedAgent
	grid         *infra.Grid
	// data recorder
	//DataRecorder *gameRecorder.ServerDataRecorder

	//server internal state
	turn int
	iteration int
	//allAgentsDead bool
	//gameRunner infra.GameRunner
	
}

type Network struct {
    Agents map[uuid.UUID]*agents.ExtendedAgent
}


func init () {
	rand.Seed(time.Now().UnixNano())
}


func (tserv *TMTServer) RunStartOfIteration(iteration int) {
	log.Printf("--------Start of iteration %v---------\n", iteration)

	//update context
	contexts := []string{"cause", "kin"} // Define possible contexts
    tserv.context = contexts[iteration%len(contexts)] // Assign context based on iteration
	tserv.grid = infra.CreateGrid(10, 10) // Create a 10x10 grid
	fmt.Printf("--------Start of iteration %d with context '%s'---------\n", iteration, tserv.context)
	//tserv.iteration = iteration
	//tserv.turn = 0


}

func (tserv *TMTServer) RunTurn(i, j int) {
	log.Printf("\n\nIteration %v, Turn %v, current agent count: %v\n", i, j, len(tserv.GetAgentMap()))
	tserv.turn = j

	// Print agent positions
    fmt.Println("Agent positions at turn:", j)
    for _, agent := range tserv.ActiveAgents {
        fmt.Printf("Agent %v at Position (%d, %d)\n", agent.NameID, agent.Position[0], agent.Position[1])
    }

	//1. Agents choose 0 or 1
	for _, agent := range tserv.ActiveAgents {
		decision := agent.DecideSacrifice()
        fmt.Printf("Agent %v made the decision: %v \n", agent.NameID, decision)
	}

	//2. Eliminate Agents
	remainingAgents := []*agents.ExtendedAgent{}
	for _, agentID := range tserv.ActiveAgents {
		if agentID.SelfSacrificeWillingness > 0.5 {
			remainingAgents = append(remainingAgents, agentID)
		} else {
			fmt.Printf("Agent %v has been eliminated\n", agentID.NameID)
		}
	}
	newActiveAgents := make(map[uuid.UUID]*agents.ExtendedAgent)
	for _, agent := range remainingAgents {
		newActiveAgents[agent.GetID()] = agent
	}
	tserv.ActiveAgents = newActiveAgents

	for _, agent := range tserv.ActiveAgents {
		agent.MoveRandomly(tserv.grid)
	}
	fmt.Printf("Turn %d: Ending with %d agents\n", tserv.turn, len(tserv.ActiveAgents))
	tserv.turn++
}

func (tserv *TMTServer) RunEndOfIteration(int) {
	log.Printf("--------End of iteration %v---------\n", tserv.iteration)
}