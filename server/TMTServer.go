package server

import (
	"fmt"
	"log"
	"maps"
	"math"
	"math/rand"

	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"gonum.org/v1/gonum/stat/distuv"

	"slices"

	agents "github.com/aaashah/TMT_Attachment/agents"
	gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)

type TMTServer struct {
	*server.BaseServer[infra.IExtendedAgent]

	Grid                         *infra.Grid
	clusterMap                   map[int][]uuid.UUID                // Map of cluster IDs to agent IDs
	ClusterEliminationData       map[int]*infra.ClusterEliminations // clusterID → ClusterEliminations
	totalRequiredEliminations    int
	totalVoluntaryEliminations   int
	lastEliminatedAgents         []infra.IExtendedAgent
	lastSelfSacrificedAgents     []infra.IExtendedAgent
	expectedChildren             float64
	neededProportionEliminations float64

	// data recorder
	//DataRecorder *gameRecorder.ServerDataRecorder
	JSONTurnLogs []gameRecorder.TurnJSONRecord
}

func CreateTMTServer(grid *infra.Grid) *TMTServer {
	tserv := &TMTServer{
		BaseServer:                   server.CreateBaseServer[infra.IExtendedAgent](10, 10, 50*time.Millisecond, 0),
		Grid:                         grid,
		clusterMap:                   make(map[int][]uuid.UUID),
		ClusterEliminationData:       make(map[int]*infra.ClusterEliminations),
		totalRequiredEliminations:    0,
		totalVoluntaryEliminations:   0,
		lastEliminatedAgents:         make([]infra.IExtendedAgent, 0),
		lastSelfSacrificedAgents:     make([]infra.IExtendedAgent, 0),
		expectedChildren:             1.9,
		neededProportionEliminations: 0.2,
		//DataRecorder: gameRecorder.CreateServerDataRecorder(),
		JSONTurnLogs: make([]gameRecorder.TurnJSONRecord, 0),
	}
	return tserv
}

func (tserv *TMTServer) GetAgentByID(agentID uuid.UUID) (infra.IExtendedAgent, bool) {
	agentMap := tserv.GetAgentMap()
	agent, exists := agentMap[agentID]
	return agent, exists
}

// Moved to TMTServer to avoid import cycle
func (tserv *TMTServer) UpdateAgentRelationship(agentAID, agentBID uuid.UUID, change float32) {
	agentA, existsA := tserv.GetAgentByID(agentAID)
	agentB, existsB := tserv.GetAgentByID(agentBID)

	if existsA && existsB {
		agentA.UpdateSocialNetwork(agentBID, change)
		agentB.UpdateSocialNetwork(agentAID, change)
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
	for i := range agentIDs {
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
	for agentID, agent := range tserv.GetAgentMap() {
		agent.UpdateRelationship(agentID, 1.0)
	}

	fmt.Printf("Social Network Initialized with %d connections.\n", edgeCount)
}

func (tserv *TMTServer) AddRelationship(agentAID, agentBID uuid.UUID, strength float32) {
	agentA, existsA := tserv.GetAgentByID(agentAID)
	agentB, existsB := tserv.GetAgentByID(agentBID)

	if existsA && existsB {
		agentA.UpdateRelationship(agentBID, strength)
		agentB.UpdateRelationship(agentAID, strength)
		//fmt.Printf("✅ Relationship established: %v ↔ %v (strength=%.2f)\n", agentAID, agentBID, strength)
	}

	// fmt.Printf("✅ Relationship established: %v ↔ %v (strength=%.2f)\n", agentAID, agentBID, strength)
}

func (tserv *TMTServer) RemoveRelationship(agentAID, agentBID uuid.UUID) {
	agentA, okA := tserv.GetAgentByID(agentAID)
	agentB, okB := tserv.GetAgentByID(agentBID)

	if okA && okB {
		agentA.RemoveRelationship(agentBID)
		agentB.RemoveRelationship(agentAID)
		// fmt.Printf("Relationship removed: %v ↔ %v\n", agentAID, agentBID)
	}
}

func (tserv *TMTServer) RunStartOfIteration(iteration int) {
	log.Printf("--------Start of iteration %v---------\n", iteration)
	fmt.Printf("--------Start of iteration %d---------\n", iteration)
	fmt.Printf("Total agents: %d\n", len(tserv.GetAgentMap()))
	// Ensure DataRecorder starts recording a new iteration
	//tserv.DataRecorder.RecordNewIteration()
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

func (tserv *TMTServer) RunTurn(i, j int) {
	log.Printf("\n\nIteration %v, Turn %v, current agent count: %v\n", i, j, len(tserv.GetAgentMap()))
	tserv.MoveAgents()
	tserv.RecordTurnJSON(i, j)
}

func (tserv *TMTServer) RunEndOfIteration(iter int) {
	log.Printf("--------End of iteration %v---------\n", iter)
	tserv.WriteIterationJSONLog(iter)
	// 2. Apply clustering (k-means)
	tserv.ApplyClustering()

	// snapshot of netowrk pre elimination
	for _, agent := range tserv.GetAgentMap() {
		networkLength := len(agent.GetNetwork())
		agent.AppendNetworkSizeHistory(networkLength)
	}

	// 4. Check for agent elimination
	tserv.updateAgentMortality()
	// 4.1 - natural deaths (old age)
	naturalDeathReport := tserv.getNaturalEliminationReport()
	tserv.applyElimination(naturalDeathReport)

	// 4.2 - unnatural deaths (sacrifice)
	sacrificialDeathReport := tserv.getSacrificialEliminationReport()
	tserv.applyElimination(sacrificialDeathReport)

	// 4.3 - create tombstones / temples
	fullDeathReport := make(map[uuid.UUID]infra.DeathInfo, len(naturalDeathReport)+len(sacrificialDeathReport))
	maps.Copy(fullDeathReport, naturalDeathReport)
	maps.Copy(fullDeathReport, sacrificialDeathReport)
	tserv.performSacrifices(fullDeathReport)

	// 5. After eliminations for agents in each cluster:
	for _, agents := range tserv.clusterMap {
		// 5.1 - Update social network (create/ cut links)
		tserv.UpdateSocialNetwork(agents)
		// 5.2 - Apply PTS protocol
		tserv.ApplyPTS(agents)
	}

	// 6. Update agent parameters
	tserv.updateClusterEliminations(fullDeathReport)
	tserv.updateAgentYsterofimia(fullDeathReport)
	tserv.updateAgentHeroism(fullDeathReport)

	// 7. Spawn new agents
	tserv.SpawnNewAgents()

	// Age up all agents
	for _, agent := range tserv.GetAgentMap() {
		agent.IncrementAge()
		// fmt.Printf("Agent %v aged to %d\n", agent.GetID(), agent.GetAge())
	}
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

	// fmt.Println("Cluster assignments:")
	// Initialize map if not done already
	if tserv.ClusterEliminationData == nil {
		tserv.ClusterEliminationData = make(map[int]*infra.ClusterEliminations)
	}

	// Record cluster sizes and set cluster history for each agent
	for clusterID, agents := range tserv.clusterMap {
		// Initialize if this cluster hasn't been tracked before
		if _, exists := tserv.ClusterEliminationData[clusterID]; !exists {
			tserv.ClusterEliminationData[clusterID] = &infra.ClusterEliminations{}
		}
		// Record size of cluster this turn
		tserv.ClusterEliminationData[clusterID].ClusterSizes = append(
			tserv.ClusterEliminationData[clusterID].ClusterSizes,
			len(agents),
		)

		// For each agent in this cluster, update their history
		for _, agentID := range agents {
			if agent, ok := tserv.GetAgentByID(agentID); ok {
				agent.AppendClusterHistory(clusterID, len(agents))
			}
		}
		// fmt.Printf("Cluster %d → %d agents\n", clusterID, len(agents))
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
				//sender.SendMessage(msg, receiver.GetID())
				sender.SendSynchronousMessage(msg, receiver.GetID())
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
	dist := distuv.Poisson{
		Lambda: tserv.expectedChildren,
		Src:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	parentPool := tserv.lastEliminatedAgents
	poolSize := len(parentPool)

	fmt.Printf("Parent pool size: %d\n", poolSize)

	if poolSize < 2 {
		fmt.Println("Not enough parents available to spawn new agents.")
		return
	}

	rand.Shuffle(poolSize, func(i, j int) {
		parentPool[i], parentPool[j] = parentPool[j], parentPool[i]
	})

	newKids := 0

	for i := 1; i < poolSize; i += 2 {
		parent1 := parentPool[i-1]
		parent2 := parentPool[i]
		childrenToSpawn := int(dist.Rand())
		newKids += childrenToSpawn
		for range childrenToSpawn {
			tserv.SpawnChild(parent1, parent2)
		}
		// fmt.Printf("Spawned %d children from %v and %v\n", childrenToSpawn, parent1.GetID(), parent2.GetID())
	}

	fmt.Printf("SPAWNED %d NEW CHILDREN\n", newKids)
}

func (tserv *TMTServer) SpawnChild(parent1, parent2 infra.IExtendedAgent) {
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

	parent1.AddDescendant(newAgent.GetID())
	// fmt.Printf("Agent type: %T\n", parent1)
	parent2.AddDescendant(newAgent.GetID())

	//add new agent to server
	tserv.AddAgent(newAgent)
	//fmt.Printf("New agent %v created from %v and %v with worldview %b\n", newAgent.GetID(), parent1.GetID(), parent2.GetID(), newWorldview)

	// add relationships in social network
	tserv.AddRelationship(parent1.GetID(), newAgent.GetID(), 0.5)
	tserv.AddRelationship(parent2.GetID(), newAgent.GetID(), 0.5)
}

func (tserv *TMTServer) UpdateProbabilityOfChildren() {
	roundEliminations := len(tserv.lastSelfSacrificedAgents)
	totalAgents := len(tserv.GetAgentMap())
	proportionOfEliminations := float64(roundEliminations) / float64(totalAgents)
	alpha := 0.2
	beta := 0.1

	if proportionOfEliminations >= tserv.neededProportionEliminations {
		tserv.expectedChildren = math.Min(tserv.expectedChildren+alpha*(1-tserv.expectedChildren), 2.5)
	} else {
		tserv.expectedChildren = math.Max(tserv.expectedChildren-beta*tserv.expectedChildren, 1.8)
	}
}

func (tserv *TMTServer) MixWorldviews(wv1, wv2 uint32) uint32 {
	mask := rand.Uint32() // or rand.Uint32() for full 32-bit mask
	return (wv1 & mask) | (wv2 &^ mask)
}

// ---------------------- Recording Turn Data ----------------------

func (tserv *TMTServer) RecordTurnJSON(iter, turn int) {
	var allAgentRecords []gameRecorder.JSONAgentRecord
	for _, agent := range tserv.GetAgentMap() {
		record := agent.RecordAgentJSON(agent)
		record.IsAlive = true
		allAgentRecords = append(allAgentRecords, record)
	}

	tombstonePositions := make([]gameRecorder.Position, len(tserv.Grid.Tombstones))
	for i, pos := range tserv.Grid.Tombstones {
		tombstonePositions[i] = gameRecorder.Position{X: pos.X, Y: pos.Y}
	}

	templePositions := make([]gameRecorder.Position, len(tserv.Grid.Temples))
	for i, pos := range tserv.Grid.Temples {
		templePositions[i] = gameRecorder.Position{X: pos.X, Y: pos.Y}
	}

	jsonLog := gameRecorder.TurnJSONRecord{
		Iteration:            iter,
		Turn:                 turn,
		Agents:               allAgentRecords,
		NumberOfAgents:       len(tserv.GetAgentMap()),
		EliminatedAgents:     AgentsToStrings(tserv.lastEliminatedAgents),
		SelfSacrificedAgents: AgentsToStrings(tserv.lastSelfSacrificedAgents),
		TombstoneLocations:   tombstonePositions,
		TempleLocations:      templePositions,
	}

	tserv.JSONTurnLogs = append(tserv.JSONTurnLogs, jsonLog)
}

func (tserv *TMTServer) WriteIterationJSONLog(iter int) {
	log := gameRecorder.IterationJSONRecord{
		Iteration: iter,
		Turns:     tserv.JSONTurnLogs,
	}

	err := gameRecorder.WriteIterationJSONLog("JSONlogs", log)
	if err != nil {
		fmt.Printf("Error writing iteration log: %v\n", err)
	}

	// Clear memory for next iteration
	tserv.JSONTurnLogs = nil
}

func AgentsToStrings(agents []infra.IExtendedAgent) []string {
	result := make([]string, len(agents))
	for i, agent := range agents {
		result[i] = agent.GetID().String()
	}
	return result
}
