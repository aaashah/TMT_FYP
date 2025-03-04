package server

import (
	"fmt"
	"log"
	"math/rand"

	//"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"

	agents "github.com/aaashah/TMT_Attachment/agents"
	gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
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
	DataRecorder *gameRecorder.ServerDataRecorder

	//server internal state
	turn int
	iteration int
	//allAgentsDead bool
	//gameRunner infra.GameRunner
	
}

// type Network struct {
//     Agents map[uuid.UUID]*agents.ExtendedAgent
// }


func init () {
	rand.Seed(time.Now().UnixNano())
}


func (tserv *TMTServer) RunStartOfIteration(iteration int) {
	log.Printf("--------Start of iteration %v---------\n", iteration)

	//update context
	contexts := []string{"cause", "kin"} // Define possible contexts
    tserv.context = contexts[iteration%len(contexts)] // Assign context based on iteration
	tserv.grid = infra.CreateGrid(10, 10) // Create a 10x10 grid
	tserv.iteration = iteration
	tserv.turn = 0
	fmt.Printf("--------Start of iteration %d with context '%s'---------\n", iteration, tserv.context)
	// Ensure DataRecorder starts recording a new iteration
	tserv.DataRecorder.RecordNewIteration()
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
	for _, agent := range tserv.ActiveAgents {
		if agent.SelfSacrificeWillingness > 0.5 {
			fmt.Printf("Agent %v has been eliminated (self-sacrificed)\n", agent.NameID)
		} else {
			remainingAgents = append(remainingAgents, agent)
		}
	}

	// Update ActiveAgents after elimination
	newActiveAgents := make(map[uuid.UUID]*agents.ExtendedAgent)
	for _, agent := range remainingAgents {
		newActiveAgents[agent.GetID()] = agent
	}
	tserv.ActiveAgents = newActiveAgents

	// 3. Move agents randomly
	for _, agent := range tserv.ActiveAgents {
		agent.MoveRandomly(tserv.grid)
	}

	fmt.Printf("Turn %d: Ending with %d agents\n", tserv.turn, len(tserv.ActiveAgents))

	// **Record turn data**
	tserv.RecordTurnInfo()
	tserv.turn++
}

func (tserv *TMTServer) RunEndOfIteration(int) {
	log.Printf("--------End of iteration %v---------\n", tserv.iteration)
}

// ---------------------- Recording Turn Data ----------------------
func (tserv *TMTServer) RecordTurnInfo() {
	// agent information
	agentRecords := []gameRecorder.AgentRecord{}

	// Log all alive agents using RecordAgentStatus
	for _, agent := range tserv.ActiveAgents {
		newAgentRecord := agent.RecordAgentStatus(agent)
		newAgentRecord.IsAlive = true
		newAgentRecord.TurnNumber = tserv.turn
		newAgentRecord.IterationNumber = tserv.iteration
		agentRecords = append(agentRecords, newAgentRecord)
	}

	// Log eliminated (sacrificed) agents
	for _, agent := range tserv.GetAgentMap() {
		if _, alive := tserv.ActiveAgents[agent.GetID()]; !alive { // If not in ActiveAgents, it's dead
			newAgentRecord := agent.RecordAgentStatus(agent)
			newAgentRecord.IsAlive = false
			newAgentRecord.TurnNumber = tserv.turn
			newAgentRecord.IterationNumber = tserv.iteration
			newAgentRecord.SpecialNote = "Eliminated" // Update special note
			agentRecords = append(agentRecords, newAgentRecord)
		}
	}

	// common information
	newInfraRecord := gameRecorder.NewInfraRecord(tserv.turn, tserv.iteration)

	tserv.DataRecorder.RecordNewTurn(agentRecords, newInfraRecord)
}