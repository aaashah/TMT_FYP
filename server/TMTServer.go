package server

import (
	"fmt"
	"maps"
	"math"
	"math/rand"

	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"

	"slices"

	"github.com/MattSScott/TMT_SOMAS/config"
	"github.com/MattSScott/TMT_SOMAS/gameRecorder"
	"github.com/MattSScott/TMT_SOMAS/infra"
	"github.com/google/uuid"
)

type TMTServer struct {
	*server.BaseServer[infra.IExtendedAgent]
	config                   config.Config
	grid                     *infra.Grid
	clusterMap               map[int][]uuid.UUID                // Map of cluster IDs to agent IDs
	clusterEliminationData   map[int]*infra.ClusterEliminations // clusterID → ClusterEliminations
	lastEliminatedAgents     []infra.IExtendedAgent
	lastSelfSacrificedAgents []infra.IExtendedAgent
	numVolunteeredAgents     int
	expectedChildren         float64
	gameRecorder             *gameRecorder.GameJSONRecord
	JSONTurnLogs             []gameRecorder.TurnJSONRecord
}

func CreateTMTServer(config config.Config) *TMTServer {
	return &TMTServer{
		BaseServer:               server.CreateBaseServer[infra.IExtendedAgent](config.NumIterations, config.NumTurns, 50*time.Millisecond, 100),
		config:                   config,
		grid:                     infra.NewGrid(infra.GRID_WIDTH, infra.GRID_HEIGHT),
		clusterMap:               make(map[int][]uuid.UUID),
		clusterEliminationData:   make(map[int]*infra.ClusterEliminations),
		lastEliminatedAgents:     make([]infra.IExtendedAgent, 0),
		lastSelfSacrificedAgents: make([]infra.IExtendedAgent, 0),
		numVolunteeredAgents:     0,
		expectedChildren:         config.InitialExpectedChildren,
		gameRecorder:             gameRecorder.MakeGameRecord(config),
		JSONTurnLogs:             make([]gameRecorder.TurnJSONRecord, 0),
	}
}

func (tserv *TMTServer) Start() {
	tserv.BaseServer.Start()
	gameRecorder.WriteJSONLog("JSONlogs", tserv.gameRecorder)
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
func (tserv *TMTServer) InitialiseRandomNetwork(p float64) {
	agentIDs := make([]uuid.UUID, 0, len(tserv.GetAgentMap()))

	// Collect all agent IDs
	for id := range tserv.GetAgentMap() {
		agentIDs = append(agentIDs, id)
	}

	if tserv.config.Debug {
		fmt.Printf("Initializing Erdős-Rényi (ER) Network with p = %.2f\n", p)
	}

	edgeCount := 0
	for i := range agentIDs {
		for j := i + 1; j < len(agentIDs); j++ { // Avoid duplicate edges
			probability := rand.Float64() // Generate a random number
			if probability <= p {         // Connect with probability p
				// Assign a random relationship strength (0.2 to 1.0)
				strength := 0.2 + rand.Float32()*0.8
				tserv.AddRelationship(agentIDs[i], agentIDs[j], strength)
				// Log connections
				if tserv.config.Debug {
					fmt.Printf("Connected Agent %v ↔ Agent %v (strength=%.2f)\n", agentIDs[i], agentIDs[j], strength)
				}
				edgeCount++
			}
		}
	}

	// add self to network
	for agentID, agent := range tserv.GetAgentMap() {
		agent.UpdateRelationship(agentID, 1.0)
	}

	if tserv.config.Debug {
		fmt.Printf("Social Network Initialized with %d connections\n", edgeCount)
	}
}

func (tserv *TMTServer) AddRelationship(agentAID, agentBID uuid.UUID, strength float32) {
	agentA, existsA := tserv.GetAgentByID(agentAID)
	agentB, existsB := tserv.GetAgentByID(agentBID)

	if existsA && existsB {
		agentA.UpdateRelationship(agentBID, strength)
		agentB.UpdateRelationship(agentAID, strength)
		if tserv.config.Debug {
			fmt.Printf("✅ Relationship established: %v ↔ %v (strength=%.2f)\n", agentAID, agentBID, strength)
		}
	}
}

func (tserv *TMTServer) RemoveRelationship(agentAID, agentBID uuid.UUID) {
	agentA, okA := tserv.GetAgentByID(agentAID)
	agentB, okB := tserv.GetAgentByID(agentBID)

	if okA && okB {
		agentA.RemoveRelationship(agentBID)
		agentB.RemoveRelationship(agentAID)
		if tserv.config.Debug {
			fmt.Printf("Relationship removed: %v ↔ %v\n", agentAID, agentBID)
		}
	}
}

func (tserv *TMTServer) RunStartOfIteration(iteration int) {
	if tserv.config.Debug {
		fmt.Printf("--------Start of iteration %d---------\n", iteration)
		fmt.Printf("Total agents: %d\n", len(tserv.GetAgentMap()))
	}
	// Clear memory for iteration
	tserv.JSONTurnLogs = nil
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
	grid := tServ.grid
	if moveX < 0 || moveX >= grid.Width {
		return false
	}
	if moveY < 0 || moveY >= grid.Height {
		return false
	}
	return !grid.IsOccupied(moveX, moveY)
}

func (tserv *TMTServer) RunTurn(i, j int) {
	if tserv.config.Debug {
		fmt.Printf("Iteration %d, Turn %d\n", i, j)
	}
	tserv.moveAgents()
	tserv.recordTurnJSON(j)
}

func (tserv *TMTServer) RunEndOfIteration(iter int) {
	if tserv.config.Debug {
		fmt.Printf("--------End of iteration %v---------\n", iter)
	}
	initialPop := len(tserv.GetAgentMap())
	// fmt.Println(len(tserv.GetAgentMap()))
	tserv.addIterationJSON(iter)
	// 2. Apply clustering (k-means)
	tserv.applyClustering()

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
		tserv.updateSocialNetwork(agents)
		// 5.2 - Apply PTS protocol
		tserv.applyPTS(agents)
	}

	// 6. Update agent parameters
	tserv.updateClusterEliminations(fullDeathReport)
	tserv.updateAgentYsterofimia(fullDeathReport)
	tserv.updateAgentHeroism(fullDeathReport)

	// 7. Spawn new agents
	tserv.updateProbabilityOfChildren()
	tserv.spawnNewAgents()

	newPop := len(tserv.GetAgentMap())
	tserv.updateAgentWorldviews(initialPop, newPop)

	// Age up all agents
	for _, agent := range tserv.GetAgentMap() {
		agent.IncrementAge()
	}

	fmt.Println()
}

// ---------------------- Helper Functions ----------------------
func runKMeans(data [][]float64, k int) []int {
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

func (tserv *TMTServer) moveAgents() {
	for _, agent := range tserv.GetAgentMap() {
		agentPos := agent.GetPosition()
		moveX, moveY := tserv.grid.GetValidMove(agentPos.X, agentPos.Y)
		targetPos, posExists := agent.GetTargetPosition(tserv.grid)

		if posExists {
			attemptX := agentPos.X - getStep(agentPos.X, targetPos.X)
			attemptY := agentPos.Y - getStep(agentPos.Y, targetPos.Y)
			if tserv.moveIsValid(attemptX, attemptY) {
				moveX, moveY = attemptX, attemptY
			}
		}

		newPos := infra.PositionVector{X: moveX, Y: moveY}
		tserv.grid.UpdateAgentPosition(agent, newPos)
		agent.SetPosition(newPos)
	}
}

func (tserv *TMTServer) applyClustering() {
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
	clusters := runKMeans(positions, k)

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

	// Initialize map if not done already
	if tserv.clusterEliminationData == nil {
		tserv.clusterEliminationData = make(map[int]*infra.ClusterEliminations)
	}

	// Record cluster sizes and set cluster history for each agent
	for clusterID, agents := range tserv.clusterMap {
		// Initialize if this cluster hasn't been tracked before
		if _, exists := tserv.clusterEliminationData[clusterID]; !exists {
			tserv.clusterEliminationData[clusterID] = &infra.ClusterEliminations{}
		}
		// Record size of cluster this turn
		tserv.clusterEliminationData[clusterID].ClusterSizes = append(
			tserv.clusterEliminationData[clusterID].ClusterSizes,
			len(agents),
		)

		// For each agent in this cluster, update their history
		for _, agentID := range agents {
			if agent, ok := tserv.GetAgentByID(agentID); ok {
				agent.AppendClusterHistory(clusterID, len(agents))
			}
		}
		if tserv.config.Debug {
			fmt.Printf("Cluster %d → %d agents\n", clusterID, len(agents))
		}
	}
}

func (tserv *TMTServer) updateSocialNetwork(cluster []uuid.UUID) {
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

func (tserv *TMTServer) applyPTS(cluster []uuid.UUID) {
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
				sender.SendSynchronousMessage(msg, receiver.GetID())
				// mark agent as getting checked on
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

func (tserv *TMTServer) updateAgentWorldviews(initialPop, newPop int) {
	trendChange := float64(newPop) / float64(tserv.config.NumAgents)
	seasonalChange := newPop - initialPop
	for _, agent := range tserv.GetAgentMap() {
		agent.UpdateWorldview(trendChange, seasonalChange)
	}
}

func (tserv *TMTServer) updateProbabilityOfChildren() {
	roundEliminations := len(tserv.lastSelfSacrificedAgents)
	totalAgents := len(tserv.GetAgentMap())
	proportionOfEliminations := float64(roundEliminations) / float64(totalAgents)

	if proportionOfEliminations >= tserv.config.PopulationRho {
		tserv.expectedChildren = math.Min(tserv.expectedChildren*1.1, tserv.config.MaxExpectedChildren)
	} else {
		tserv.expectedChildren = math.Max(tserv.expectedChildren*0.9, tserv.config.MinExpectedChildren)
	}
}

// func (tserv *TMTServer) mixWorldviews(wv1, wv2 uint32) uint32 {
// 	mask := rand.Uint32()
// 	return (wv1 & mask) | (wv2 &^ mask)
// }

// ---------------------- Recording Turn Data ----------------------

func (tserv *TMTServer) recordTurnJSON(turn int) {
	var allAgentRecords []gameRecorder.JSONAgentRecord
	for _, agent := range tserv.GetAgentMap() {
		record := agent.RecordAgentJSON(agent)
		record.IsAlive = true
		allAgentRecords = append(allAgentRecords, record)
	}

	tombstonePositions := make([]gameRecorder.Position, len(tserv.grid.Tombstones))
	for i, pos := range tserv.grid.Tombstones {
		tombstonePositions[i] = gameRecorder.Position{X: pos.X, Y: pos.Y}
	}

	templePositions := make([]gameRecorder.Position, len(tserv.grid.Temples))
	for i, pos := range tserv.grid.Temples {
		templePositions[i] = gameRecorder.Position{X: pos.X, Y: pos.Y}
	}

	totalAgents := float64(len(tserv.GetAgentMap()))
	reqElims := int(tserv.config.PopulationRho * totalAgents)

	jsonLog := gameRecorder.TurnJSONRecord{
		Turn:                      turn,
		Agents:                    allAgentRecords,
		NumberOfAgents:            len(tserv.GetAgentMap()),
		EliminatedAgents:          agentsToStrings(tserv.lastEliminatedAgents),
		TotalRequiredEliminations: reqElims,
		TotalVolunteers:           tserv.numVolunteeredAgents,
		SelfSacrificedAgents:      agentsToStrings(tserv.lastSelfSacrificedAgents),
		TombstoneLocations:        tombstonePositions,
		TempleLocations:           templePositions,
	}

	tserv.JSONTurnLogs = append(tserv.JSONTurnLogs, jsonLog)
}

func (tserv *TMTServer) addIterationJSON(iter int) {
	log := gameRecorder.IterationJSONRecord{
		Iteration: iter,
		Turns:     tserv.JSONTurnLogs,
	}

	tserv.gameRecorder.AddIteration(log)
}

func agentsToStrings(agents []infra.IExtendedAgent) []string {
	result := make([]string, len(agents))
	for i, agent := range agents {
		result[i] = agent.GetID().String()
	}
	return result
}
