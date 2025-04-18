package server

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	//"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"

	"slices"

	agents "github.com/aaashah/TMT_Attachment/agents"
	gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)

type TMTServer struct {
	*server.BaseServer[infra.IExtendedAgent]

	//agentInfoList []infra.IExtendedAgent
	//mu     sync.Mutex

	Grid        *infra.Grid
	PositionMap map[[2]int]*agents.ExtendedAgent // Map of agent positions
	clusterMap  map[int][]uuid.UUID              // Map of cluster IDs to agent IDs

	// data recorder
	DataRecorder *gameRecorder.ServerDataRecorder

	//server internal state
	turn      int
	iteration int
	//allAgentsDead bool
	//gameRunner infra.GameRunner

}

func init() {
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

	if !existsA || !existsB {
		return
	}

	agentA.UpdateSocialNetwork(agentBID, change)

	agentB.UpdateSocialNetwork(agentAID, change)

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
				//agentA := tserv.GetAgentMap()[agentIDs[i]]
				//agentB := tserv.GetAgentMap()[agentIDs[j]]

				// Assign a random relationship strength (0.2 to 1.0)
				strength := 0.2 + rand.Float32()*0.8
				//tserv.AddRelationship(agentA.GetID(), agentB.GetID(), strength)
				tserv.AddRelationship(agentIDs[i], agentIDs[j], strength)

				// Log connections
				//fmt.Printf("Connected Agent %v ↔ Agent %v (strength=%.2f)\n",
				//agentA.GetID(), agentB.GetID(), strength)
				edgeCount++
			}
		}
	}

	// add self to network
	for _, agent := range tserv.GetAgentMap() {
		agentID := agent.GetID()
		agent.UpdateRelationship(agentID, 1.0) 
	}

	fmt.Printf("Social Network Initialized with %d connections.\n", edgeCount)
}

func (tserv *TMTServer) AddRelationship(agentAID, agentBID uuid.UUID, strength float32) {
	agentA, existsA := tserv.GetAgentMap()[agentAID]
	agentB, existsB := tserv.GetAgentMap()[agentBID]

	if existsA && existsB {
		agentA.UpdateRelationship(agentBID, strength)
		agentB.UpdateRelationship(agentAID, strength)
		//fmt.Printf("✅ Relationship established: %v ↔ %v (strength=%.2f)\n", agentAID, agentBID, strength)
	}

	fmt.Printf("✅ Relationship established: %v ↔ %v (strength=%.2f)\n", agentAID, agentBID, strength)
}

func (tserv *TMTServer) RemoveRelationship(agentAID, agentBID uuid.UUID) {
	agentA, okA := tserv.GetAgentByID(agentAID)
	agentB, okB := tserv.GetAgentByID(agentBID)

	if okA && okB {
		agentA.RemoveRelationship(agentBID)
		agentB.RemoveRelationship(agentAID)

		fmt.Printf("Relationship removed: %v ↔ %v\n", agentAID, agentBID)
	}
}

func (tserv *TMTServer) RunStartOfIteration(iteration int) {
	log.Printf("--------Start of iteration %v---------\n", iteration)

	tserv.iteration = iteration
	tserv.turn = 0

	if iteration == 0 {
		const connectionProbability = 0.35
		tserv.InitialiseRandomNetwork(connectionProbability)
	}

	// Age up all agents at the start of each iteration (except start of game)
	// if iteration > 0 {
	for _, agent := range tserv.GetAgentMap() {
		agent.IncrementAge()
		fmt.Printf("Agent %v aged to %d\n", agent.GetID(), agent.GetAge())
	}
	// }


	fmt.Printf("--------Start of iteration %d---------\n", iteration)
	// Ensure DataRecorder starts recording a new iteration
	tserv.DataRecorder.RecordNewIteration()
}

func getStep(current, target int) int {
	if target > current {
		return 1
	} else if target < current {
		return -1
	}
	return 0
}

func (tServ *TMTServer) moveIsValid(moveX, moveY int) bool {
	// moveX >= 0 && moveX < grid.Width && moveY >= 0 && moveY < grid.Height && !grid.IsOccupied(moveX, moveY)
	grid := tServ.Grid
	if moveX < 0 || moveX >= grid.Width {
		return false
	}
	if moveY < 0 || moveY >= grid.Height {
		return false
	}
	return !grid.IsOccupied(moveX, moveY)
}

const MOVEMENT_TURNS int = 20

func (tserv *TMTServer) RunTurn(i, j int) {
	log.Printf("\n\nIteration %v, Turn %v, current agent count: %v\n", i, j, len(tserv.GetAgentMap()))
	tserv.turn = j
	// if i == 0 && j == 0 {
	// 	tserv.RecordTurnInfo()
	// 	return
	// }

	// 1. Move agents
	for range MOVEMENT_TURNS {
		tserv.MoveAgents()
	}

	// 2. Apply clustering (k-means)
	tserv.ApplyClustering()

	// 4. Check for agent elimination
	//tserv.ApplyAPS()
	tserv.ApplyElimination(j)

	// 5. After eliminations for agents in each cluster:
	for _, agents := range tserv.clusterMap {
		// age up agents here??

		// 5.1 Update social network (create/ cut links)
		tserv.UpdateSocialNetwork(agents)
		// 5.2 apply PTS protocol
		tserv.ApplyPTS(agents)
		// 5.3 update heroism
		//(done within ApplyElimination)
		// 5.4 update worldview ??
	}


	fmt.Printf("Turn %d: Ending with %d agents\n", tserv.turn, len(tserv.GetAgentMap()))
	tserv.RecordTurnInfo()
}

func (tserv *TMTServer) RunEndOfIteration(int) {
	log.Printf("--------End of iteration %v---------\n", tserv.iteration)
	//tserv.iteration++
	// spawn new agents
	tserv.SpawnNewAgents()
}

// ---------------------- Helper Functions ----------------------
func RunKMeans(data [][]float64, k int) []int {
	if len(data) == 0 {
		return []int{}
	}
	// Initialize centroids
	centroids := make([][]float64, k)
	for i := range k {
		centroids[i] = slices.Clone(data[rand.Intn(len(data))])
	}

	assignments := make([]int, len(data))
	changed := true

	for changed {
		changed = false

		// Assign points
		for i, point := range data {
			minDist := math.MaxFloat64
			best := -1
			for j, centroid := range centroids {
				d := distance(point, centroid)
				if d < minDist {
					minDist = d
					best = j
				}
			}
			if assignments[i] != best {
				assignments[i] = best
				changed = true
			}
		}

		// Update centroids
		count := make([]int, k)
		sums := make([][]float64, k)
		for i := range sums {
			sums[i] = make([]float64, 2) // because [x, y]
		}
		for i, a := range assignments {
			sums[a][0] += data[i][0]
			sums[a][1] += data[i][1]
			count[a]++
		}
		for i := range k {
			if count[i] > 0 {
				centroids[i][0] = sums[i][0] / float64(count[i])
				centroids[i][1] = sums[i][1] / float64(count[i])
			}
		}
	}

	return assignments
}

func distance(a, b []float64) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	return math.Sqrt(dx*dx + dy*dy)
}

func (tserv *TMTServer) MoveAgents() {
	for _, agent := range tserv.GetAgentMap() {
		agentPos := agent.GetPosition()
		moveX, moveY := tserv.Grid.GetValidMove(agentPos.X, agentPos.Y)
		targetPos, posExists := agent.GetTargetPosition(tserv.Grid)

		if posExists {
			attemptX := agentPos.X - getStep(agentPos.X, targetPos.X)
			attemptY := agentPos.Y - getStep(agentPos.Y, targetPos.Y)
			if tserv.moveIsValid(attemptX, attemptY) {
				moveX, moveY = attemptX, attemptY
			}
		}

		//tserv.Grid.UpdateAgentPosition(agent, moveX, moveY)
		newPos := infra.PositionVector{X: moveX, Y: moveY}
		tserv.Grid.UpdateAgentPosition(agent, newPos)
		agent.SetPosition(newPos)
	}
}

func (tserv *TMTServer) ApplyClustering() {
	positions := [][]float64{}
	idToIndex := make([]uuid.UUID, 0)
	agentMap := tserv.GetAgentMap()
	if len(agentMap) == 0 {
		return // Nothing to cluster
	}

	for _, agent := range agentMap {
		pos := agent.GetPosition()
		positions = append(positions, []float64{float64(pos.X), float64(pos.Y)})
		idToIndex = append(idToIndex, agent.GetID())
	}

	k := 3
	clusters := RunKMeans(positions, k)

	for i, clusterID := range clusters {
		agentID := idToIndex[i]
		if agent, ok := tserv.GetAgentByID(agentID); ok {
			agent.SetClusterID(clusterID)
		}
	}
	tserv.clusterMap = make(map[int][]uuid.UUID)
	for _, agent := range tserv.GetAgentMap() {
		tserv.clusterMap[agent.GetClusterID()] = append(tserv.clusterMap[agent.GetClusterID()], agent.GetID())
	}

	fmt.Println("Cluster assignments:")
	for clusterID, agents := range tserv.clusterMap {
		fmt.Printf("Cluster %d → %d agents\n", clusterID, len(agents))
	}
}

func (tserv *TMTServer) ApplyElimination(turn int) {
	agentsToRemove := make(map[uuid.UUID]bool)

	if turn == 0 {
		tserv.updateAgentMortality()
		for _, agent := range tserv.GetAgentMap() {
			if !agent.IsAlive() {
				fmt.Printf("Agent %v has been eliminated (natural causes)\n", agent.GetID())
				pos := agent.GetPosition()
				tserv.Grid.PlaceTombstone(pos.X, pos.Y)
				agentsToRemove[agent.GetID()] = true
			}
		}
	} else {
		// voluntary self-sacrifice
		for _, agent := range tserv.GetAgentMap() {
			if agent.GetASPDecision(tserv.Grid) == infra.SELF_SACRIFICE {
				// 4.1 Place temples/monuments for self-sacrificed agents
				fmt.Printf("Agent %v has been eliminated (self-sacrificed)\n", agent.GetID())
				pos := agent.GetPosition()
				tserv.Grid.PlaceTemple(pos.X, pos.Y)
				agent.IncrementHeroism()
				agentsToRemove[agent.GetID()] = true
			}
		}
		// pick random agent to eliminate if no one self-sacrificed
		if len(agentsToRemove) == 0 {
			livingAgents := []infra.IExtendedAgent{}
			for _, agent := range tserv.GetAgentMap() {
				if agent.IsAlive() {
					livingAgents = append(livingAgents, agent)
				}
			}
			if len(livingAgents) > 0 {
				victim := livingAgents[rand.Intn(len(livingAgents))]
				fmt.Printf("No self-sacrifices this turn — randomly eliminating Agent %v\n", victim.GetID())
				pos := victim.GetPosition()
				tserv.Grid.PlaceTombstone(pos.X, pos.Y)
				agentsToRemove[victim.GetID()] = true
			}
		}
	}

	// also track eliminations per cluster and in network
	clusterEliminationCount := make(map[int]int) // number of eliminations per cluster
	for id := range agentsToRemove {
		agent, ok := tserv.GetAgentByID(id)
		if ok {
			clusterID := agent.GetClusterID()    // get the cluster ID of the agent
			clusterEliminationCount[clusterID]++ // increment the count for that cluster
			tserv.RemoveAgent(agent)
		}
	}

	// update history for remaining agents:
	for _, agent := range tserv.GetAgentMap() {
		//cluster eliminations
		clusterID := agent.GetClusterID()
		if eliminatedInCluster, exists := clusterEliminationCount[clusterID]; exists {
			agent.IncrementClusterEliminations(eliminatedInCluster)
			//fmt.Printf("Agent %v in cluster %d has %d eliminations in this cluster\n", agent.GetID(), clusterID, eliminatedInCluster)
		}

		// social network eliminations
		networkEliminationCount := 0
		for friendID, esteem := range agent.GetNetwork() {
			if agentsToRemove[friendID] {
				eliminatedAgent, exists := tserv.GetAgentByID(friendID)
				if !exists {
					continue
				}
				networkEliminationCount++
				
				if eliminatedAgent.GetASPDecision(tserv.Grid) == infra.SELF_SACRIFICE {
					ysterofimia := agent.GetYsterofimia()
					ysterofimia.IncrementSelfSacrificeCount()
					ysterofimia.AddSelfSacrificeEsteems(esteem)
				} else {
					ysterofimia := agent.GetYsterofimia()
					ysterofimia.IncrementOtherEliminationCount()
					ysterofimia.AddOtherEliminationsEsteems(esteem)
				}
			}
		}
		agent.IncrementNetworkEliminations(networkEliminationCount)

		//track eliminations and esteems for ysterofimia
	}

}

func (tserv *TMTServer) updateAgentMortality() {
	for _, agent := range tserv.GetAgentMap() {
		probDeath := agent.GetTelomere()
		randVal := rand.Float32()
		if randVal < probDeath {
			agent.MarkAsDead()
		}
	}
}

func (tserv *TMTServer) UpdateSocialNetwork(cluster []uuid.UUID) {
	for _, agentID := range cluster {
		agent, ok := tserv.GetAgentByID(agentID)
		if !ok || !agent.IsAlive() {
			continue
		}

		for _, otherID := range cluster {
			if otherID == agentID {
				continue
			}

			_, connected := agent.GetNetwork()[otherID]

			// add link based on anxiety
			if !connected && rand.Float32() < agent.GetAttachment().Anxiety {
				// Add symmetric relationship with random strength
				strength := 0.2 + rand.Float32()*0.8
				tserv.AddRelationship(agentID, otherID, strength)
			}

			// remove link based on attachment
			if connected && rand.Float32() < agent.GetAttachment().Avoidance {
				agent.RemoveRelationship(otherID)
				tserv.RemoveRelationship(agentID, otherID)
			}
		}
	}
}

func (tserv *TMTServer) ApplyPTS(cluster []uuid.UUID) {
	receivedCheck := make(map[uuid.UUID]bool) // Track who got a check
	// get ExtendedAgent from agentID
	for _, senderID := range cluster {
		sender, ok := tserv.GetAgentByID(senderID)
		if !ok {
			continue
		}
		if rand.Float32() < sender.GetPTSParams().CheckProb {
			for _, receiverID := range cluster {
				if receiverID == senderID {
					continue // don't send to self
				}
				receiver, ok := tserv.GetAgentByID(receiverID)
				if !ok {
					continue
				}
				// send wellbeing check message
				msg := sender.CreateWellbeingCheckMessage()
				sender.SendMessage(msg, receiver.GetID())
				//fmt.Printf("Agent %v sent wellbeing check to %v\n", sender.GetID(), receiverID)

				//make agent as getting checked on
				receivedCheck[receiver.GetID()] = true
			}
		}
	}

	//update beta for agents that didn't get checked on
	for _, agentID := range cluster {
		if receivedCheck[agentID] {
			continue // skip agents who received a check
		}
		agent, ok := tserv.GetAgentByID(agentID)
		if !ok {
			continue
		}

		//update beta
		for neighbourID := range agent.GetNetwork() {
			agent.UpdateEsteem(neighbourID, false)
		}
	}
}

func (tserv *TMTServer) SpawnNewAgents() {
	agentMap := tserv.GetAgentMap()

	parentIDs := make([]uuid.UUID, 0, len(agentMap))
	for id := range agentMap {
		parentIDs = append(parentIDs, id)
	}

	rand.Shuffle(len(parentIDs), func(i, j int) { parentIDs[i], parentIDs[j] = parentIDs[j], parentIDs[i] })

	for i := 0; i+1 < len(parentIDs); i += 2 {
		parent1, _ := tserv.GetAgentByID(parentIDs[i])
		parent2, _ := tserv.GetAgentByID(parentIDs[i+1])

		if !parent1.IsAlive() || !parent2.IsAlive() {
			continue
		}

		newWorldview := tserv.MixWorldviews(parent1.GetWorldviewBinary(), parent2.GetWorldviewBinary())
			randVal := rand.Float32()

			var newAgent infra.IExtendedAgent

			switch {
			case randVal < 0.25:
				newAgent = agents.CreateSecureAgent(tserv, tserv.Grid, parent1.GetID(), parent2.GetID(), newWorldview)
			case randVal < 0.5:
				newAgent = agents.CreateDismissiveAgent(tserv, tserv.Grid, parent1.GetID(), parent2.GetID(), newWorldview)
			case randVal < 0.75:
				newAgent = agents.CreatePreoccupiedAgent(tserv, tserv.Grid, parent1.GetID(), parent2.GetID(), newWorldview)
			default:
				newAgent = agents.CreateFearfulAgent(tserv, tserv.Grid, parent1.GetID(), parent2.GetID(), newWorldview)
			}

			//create new agent
			//newAgent := agents.CreateSecureAgent(tserv, tserv.Grid)

			//newAgent.SetWorldviewBinary(newWorldview)
			//newAgent.SetParents(parent1.GetID(), parent2.GetID())
			parent1.AddDescendant(newAgent.GetID())
			parent2.AddDescendant(newAgent.GetID())

			//add new agent to server
			tserv.AddAgent(newAgent)
			//fmt.Printf("New agent %v created from %v and %v with worldview %b\n", newAgent.GetID(), parent1.GetID(), parent2.GetID(), newWorldview)

			// add relationships in social network
			tserv.AddRelationship(parent1.GetID(), newAgent.GetID(), 0.5)
			tserv.AddRelationship(parent2.GetID(), newAgent.GetID(), 0.5)
	}
}

func (tserv *TMTServer) MixWorldviews(wv1, wv2 uint32) uint32 {
	mask := uint32(rand.Int31()) // or rand.Uint32() for full 32-bit mask
	return (wv1 & mask) | (wv2 &^ mask)
}

// ---------------------- Recording Turn Data ----------------------
func (tserv *TMTServer) RecordTurnInfo() {
	// Create a new infra record
	newInfraRecord := gameRecorder.NewInfraRecord(tserv.turn, tserv.iteration)

	// Record agent positions
	for _, agent := range tserv.GetAgentMap() {
		pos := agent.GetPosition()
		newInfraRecord.AgentPositions[[2]int{pos.X, pos.Y}] = true
	}

	// Record tombstone locations
	for _, tombstonePos := range tserv.Grid.Tombstones {
		newInfraRecord.Tombstones[[2]int{tombstonePos.X, tombstonePos.Y}] = true
	}

	for _, templePos := range tserv.Grid.Temples {
		newInfraRecord.Temples[[2]int{templePos.X, templePos.Y}] = true
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
