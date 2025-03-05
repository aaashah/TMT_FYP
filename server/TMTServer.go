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
	//context string
	ActiveAgents map[uuid.UUID]*agents.ExtendedAgent
	grid         *infra.Grid
	PositionMap map[[2]int]*agents.ExtendedAgent // Map of agent positions

	// data recorder
	DataRecorder *gameRecorder.ServerDataRecorder

	//server internal state
	turn int
	iteration int
	//allAgentsDead bool
	//gameRunner infra.GameRunner
	
}

var _ infra.IServer = (*TMTServer)(nil)

func init () {
	rand.Seed(time.Now().UnixNano())
}

func (tserv *TMTServer) GetAgentByID(agentID uuid.UUID) (infra.IExtendedAgent, bool) {
	agent, exists := tserv.ActiveAgents[agentID]
	return agent, exists
}

func (tserv *TMTServer) GetAgentPosition(agentID uuid.UUID) ([2]int, bool) {
	agent, exists := tserv.ActiveAgents[agentID]
	if !exists {
		return [2]int{0, 0}, false
	}
	return agent.Position, true
}

// Moved to TMTServer to avoid import cycle
func (tserv *TMTServer) UpdateAgentRelationship(agentAID, agentBID uuid.UUID, change float32) {
	agentA, existsA := tserv.ActiveAgents[agentAID]
	agentB, existsB := tserv.ActiveAgents[agentBID]

	if existsA && existsB {
		agentA.UpdateRelationship(agentBID, change)
		agentB.UpdateRelationship(agentAID, change)
	}
}

// Erdős–Rényi (ER) Random Network
func (tserv *TMTServer) InitialiseRandomNetwork(p float32) {
	agentIDs := make([]uuid.UUID, 0, len(tserv.ActiveAgents))

	// Collect all agent IDs
	for id := range tserv.ActiveAgents {
		agentIDs = append(agentIDs, id)
	}

	fmt.Printf("Initializing Erdős-Rényi (ER) Network with p = %.2f\n", p)

	edgeCount := 0
	for i := 0; i < len(agentIDs); i++ {
		for j := i + 1; j < len(agentIDs); j++ { // Avoid duplicate edges
			probability := rand.Float32() // Generate a random number
			// fmt.Printf("Checking link between %v and %v (p=%.2f, rolled=%.2f)\n", 
			// 	agentIDs[i], agentIDs[j], p, probability)
			
			if probability <= p { // Connect with probability p
				agentA := tserv.ActiveAgents[agentIDs[i]]
				agentB := tserv.ActiveAgents[agentIDs[j]]

				// Assign a random relationship strength (0.2 to 1.0)
				strength := 0.2 + rand.Float32()*0.8
				tserv.AddRelationship(agentA.GetID(), agentB.GetID(), strength)

				// Log connections
				fmt.Printf("Connected Agent %v ↔ Agent %v (strength=%.2f)\n",
					agentA.GetID(), agentB.GetID(), strength)
				edgeCount++
			}
		}
	}

	fmt.Printf("Social Network Initialized with %d connections.\n", edgeCount)
}

func (tserv *TMTServer) AddRelationship(agentAID, agentBID uuid.UUID, strength float32) {
	agentA, existsA := tserv.ActiveAgents[agentAID]
	agentB, existsB := tserv.ActiveAgents[agentBID]

	if existsA && existsB {
		agentA.Network[agentBID] = strength
		agentB.Network[agentAID] = strength
		fmt.Printf("✅ Relationship established: %v ↔ %v (strength=%.2f)\n", agentAID, agentBID, strength)
	} else {
		fmt.Printf("❌ Relationship failed: %v ↔ %v (agent missing?)\n", agentAID, agentBID)
	}
}

func (tserv *TMTServer) RunStartOfIteration(iteration int) {
	log.Printf("--------Start of iteration %v---------\n", iteration)

	tserv.grid = infra.CreateGrid(10, 10) // Create a 10x10 grid
	tserv.iteration = iteration
	tserv.turn = 0

	// Print the network structure
	fmt.Println("Agent Social Network at iteration start:")
	for _, agent := range tserv.ActiveAgents {
		fmt.Printf("Agent %v connections: ", agent.GetID())
		for otherID, strength := range agent.GetNetwork() {
			fmt.Printf("(%v, strength=%.2f) ", otherID, strength)
		}
		fmt.Println()
	}
	
	fmt.Printf("--------Start of iteration %d---------\n", iteration)
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

	// 1. Move agents based on social network attraction
	for _, agent := range tserv.ActiveAgents {
		agent.MoveAttractedToNetwork(tserv.grid, tserv)
	}

	// 2. Agents make decisions
	for _, agent := range tserv.ActiveAgents {
		decision := agent.DecideSacrifice()
		fmt.Printf("Agent %v willing to sacrifice by: %v \n", agent.NameID, decision)
	}

	//3. Eliminate Agents
	remainingAgents := []*agents.ExtendedAgent{}
	for _, agent := range tserv.ActiveAgents {
		if agent.SelfSacrificeWillingness > 0.3 {
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

	// 3. Move agents based on network
	for _, agent := range tserv.ActiveAgents {
		agent.MoveAttractedToNetwork(tserv.grid, tserv)
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

