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
	//ActiveAgents map[uuid.UUID]*agents.ExtendedAgent
	Grid         *infra.Grid
	PositionMap map[[2]int]*agents.ExtendedAgent // Map of agent positions

	// data recorder
	DataRecorder *gameRecorder.ServerDataRecorder

	//server internal state
	turn int
	iteration int
	//allAgentsDead bool
	//gameRunner infra.GameRunner
	
}


func init () {
	rand.Seed(time.Now().UnixNano())
}

func (tserv *TMTServer) GetAgentByID(agentID uuid.UUID) (infra.IExtendedAgent, bool) {
	agent, exists := tserv.GetAgentMap()[agentID]
	return agent, exists
}



// Moved to TMTServer to avoid import cycle
func (tserv *TMTServer) UpdateAgentRelationship(agentAID, agentBID uuid.UUID, change float32) {
	agentA, existsA := tserv.GetAgentMap()[agentAID]
	agentB, existsB := tserv.GetAgentMap()[agentBID]

	if existsA && existsB {
		agentA.UpdateRelationship(agentBID, change)
		agentB.UpdateRelationship(agentAID, change)
	}
}

// Erdős–Rényi (ER) Random Network
func (tserv *TMTServer) InitialiseRandomNetwork(p float32) {
	agentIDs := make([]uuid.UUID, 0, len(tserv.GetAgentMap()))

	// Collect all agent IDs
	for id := range tserv.GetAgentMap() {
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
				agentA := tserv.GetAgentMap()[agentIDs[i]]
				agentB := tserv.GetAgentMap()[agentIDs[j]]

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
	agentA, existsA := tserv.GetAgentMap()[agentAID]
	agentB, existsB := tserv.GetAgentMap()[agentBID]

	if existsA && existsB {
		agentA.GetNetwork()[agentBID] = strength
		agentB.GetNetwork()[agentAID] = strength
		fmt.Printf("✅ Relationship established: %v ↔ %v (strength=%.2f)\n", agentAID, agentBID, strength)
	} else {
		fmt.Printf("❌ Relationship failed: %v ↔ %v (agent missing?)\n", agentAID, agentBID)
	}
}

func (tserv *TMTServer) RunStartOfIteration(iteration int) {
	log.Printf("--------Start of iteration %v---------\n", iteration)

	tserv.iteration = iteration
	tserv.turn = 0

	// Age up all agents at the start of each iteration (except start of game)
	if iteration > 0 {
		for _, agent := range tserv.GetAgentMap() {
			agent.SetAge(agent.GetAge() + 1)
			fmt.Printf("Agent %v aged to %d\n", agent.GetID(), agent.GetAge())
		}
	}
	
	// Print the network structure
	fmt.Println("Agent Social Network at iteration start:")
	for _, agent := range tserv.GetAgentMap() {
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

	for _, agent := range tserv.GetAgentMap() {
		agent.Move(tserv.Grid)
	}

	if i == 0 && j == 0 {
		tserv.RecordTurnInfo()
		tserv.turn++
		return
	}

	for _, agent := range tserv.GetAgentMap() {
		decision := agent.DecideSacrifice()
		fmt.Printf("Agent %v willing to sacrifice by: %v \n", agent.GetID(), decision)
	}

	agentsToRemove := make(map[uuid.UUID]bool)

	if j == 0 {
		for _, agent := range tserv.GetAgentMap() {
			if agent.GetMortality() {
				fmt.Printf("Agent %v has been eliminated (natural causes)\n", agent.GetID())
				pos := agent.GetPosition()
				tserv.Grid.PlaceTombstone(pos[0], pos[1])
				agentsToRemove[agent.GetID()] = true
			}
		}
	} else {
		for _, agent := range tserv.GetAgentMap() {
			if agent.GetSelfSacrificeWillingness() > 0.9 {
				fmt.Printf("Agent %v has been eliminated (self-sacrificed)\n", agent.GetID())
				pos := agent.GetPosition()
				tserv.Grid.PlaceTemple(pos[0], pos[1])
				agentsToRemove[agent.GetID()] = true
			}
		}
	}

	for id := range agentsToRemove {
		agent, ok := tserv.GetAgentByID(id)
		if ok {
			tserv.RemoveAgent(agent)
		}
	}

	fmt.Printf("Turn %d: Ending with %d agents\n", tserv.turn, len(tserv.GetAgentMap()))
	tserv.RecordTurnInfo()
	tserv.turn++
}

func (tserv *TMTServer) RunEndOfIteration(int) {
	log.Printf("--------End of iteration %v---------\n", tserv.iteration)
}

// ---------------------- Recording Turn Data ----------------------
func (tserv *TMTServer) RecordTurnInfo() {
	// Create a new infra record
	newInfraRecord := gameRecorder.NewInfraRecord(tserv.turn, tserv.iteration)

	// Record agent positions
	for _, agent := range tserv.GetAgentMap() {
		pos := agent.GetPosition()
		newInfraRecord.AgentPositions[[2]int{pos[0], pos[1]}] = true
	}

	// Record tombstone locations
	for tombstonePos := range tserv.Grid.Tombstones {
		newInfraRecord.Tombstones[tombstonePos] = true
	}

	for templePos := range tserv.Grid.Temples {
		newInfraRecord.Temples[templePos] = true
	}

	// Collect agent records
	agentRecords := []gameRecorder.AgentRecord{}
	for _, agent := range tserv.GetAgentMap() {
		newAgentRecord := agent.RecordAgentStatus(agent)
		newAgentRecord.IsAlive = true
		newAgentRecord.TurnNumber = tserv.turn
		newAgentRecord.IterationNumber = tserv.iteration
		// Explicitly fetch the latest age instead of using stale data
		newAgentRecord.AgentAge = agent.GetAge()
		//fmt.Printf("[DEBUG] Recorded Age for Agent %v: %d\n", agent.GetID(), newAgentRecord.AgentAge)
		agentRecords = append(agentRecords, newAgentRecord)
	}

	// Record eliminated agents
	for _, agent := range tserv.GetAgentMap() {
		if _, alive := tserv.GetAgentMap()[agent.GetID()]; !alive { 
			newAgentRecord := agent.RecordAgentStatus(agent)
			newAgentRecord.IsAlive = false
			newAgentRecord.TurnNumber = tserv.turn
			newAgentRecord.IterationNumber = tserv.iteration
			//newAgentRecord.Died = agent.GetMortality()
			newAgentRecord.SpecialNote = "Eliminated"

			// Explicitly store the last known age before elimination
			newAgentRecord.AgentAge = agent.GetAge()
			//fmt.Printf("[DEBUG] Recorded Age for Agent %v: %d\n", agent.GetID(), newAgentRecord.AgentAge)
			agentRecords = append(agentRecords, newAgentRecord)
		}
	}

	// Save the recorded turn in the data recorder
	tserv.DataRecorder.RecordNewTurn(agentRecords, newInfraRecord)
}
