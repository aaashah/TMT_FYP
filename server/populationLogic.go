package server

import (
	"fmt"
	"math/rand"

	"maps"

	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)

func (tserv *TMTServer) updateAgentMortality() {
	for _, agent := range tserv.GetAgentMap() {
		probDeath := agent.GetTelomere()
		randVal := rand.Float32()
		if randVal < probDeath {
			agent.MarkAsDead()
		}
	}
}

func (tserv *TMTServer) getNaturalEliminations() map[uuid.UUID]infra.IExtendedAgent {
	naturalElims := make(map[uuid.UUID]infra.IExtendedAgent)
	for agentID, agent := range tserv.GetAgentMap() {
		if !agent.IsAlive() {
			fmt.Printf("Agent %v has been eliminated (natural causes)\n", agent.GetID())
			pos := agent.GetPosition()
			tserv.Grid.PlaceTombstone(pos.X, pos.Y)
			// naturalElims = append(naturalElims, agent)
			naturalElims[agentID] = agent
			tserv.LastEliminatedAgents = append(tserv.LastEliminatedAgents, agent)
		}
	}
	return naturalElims
}

func (tserv *TMTServer) stratifyVolunteers() ([]infra.IExtendedAgent, []infra.IExtendedAgent) {
	volunteers := make([]infra.IExtendedAgent, 0)
	nonVolunteers := make([]infra.IExtendedAgent, 0)

	// Separate volunteers and non-volunteers
	for _, agent := range tserv.GetAgentMap() {
		// check if agent volunteered
		if agent.GetASPDecision(tserv.Grid) == infra.SELF_SACRIFICE {
			volunteers = append(volunteers, agent)
		} else {
			nonVolunteers = append(nonVolunteers, agent)
		}
	}
	return volunteers, nonVolunteers
}

func (tserv *TMTServer) voluntarilySacrificeAgent(agent infra.IExtendedAgent) {
	pos := agent.GetPosition()
	tserv.Grid.PlaceTemple(pos.X, pos.Y)
	agent.IncrementHeroism()
	tserv.LastEliminatedAgents = append(tserv.LastEliminatedAgents, agent)
	tserv.LastSelfSacrificedAgents = append(tserv.LastSelfSacrificedAgents, agent)
	fmt.Printf("Agent %v has been eliminated (voluntary)\n", agent.GetID())
}

func (tserv *TMTServer) involuntarilySacrificeAgent(agent infra.IExtendedAgent) {
	pos := agent.GetPosition()
	tserv.Grid.PlaceTombstone(pos.X, pos.Y)
	tserv.LastEliminatedAgents = append(tserv.LastEliminatedAgents, agent)
	fmt.Printf("Agent %v has been eliminated (non-voluntary)\n", agent.GetID())
}

func (tserv *TMTServer) getSacrificialEliminations(volunteers, nonVolunteers []infra.IExtendedAgent) map[uuid.UUID]infra.IExtendedAgent {
	sacrificialElims := make(map[uuid.UUID]infra.IExtendedAgent)
	totalAgents := float64(len(tserv.GetAgentMap()))
	neededVolunteers := int(tserv.neededProportionEliminations * totalAgents)
	actualVolunteers := len(volunteers)

	if actualVolunteers >= neededVolunteers {
		//randomly select n volunteers to eliminate
		rand.Shuffle(actualVolunteers, func(i, j int) {
			volunteers[i], volunteers[j] = volunteers[j], volunteers[i]
		})
		for i := range neededVolunteers {
			agent := volunteers[i]
			agentID := agent.GetID()
			tserv.voluntarilySacrificeAgent(agent)
			sacrificialElims[agentID] = agent
		}
	} else {
		//eliminate all volunteers...
		for _, agent := range volunteers {
			agentID := agent.GetID()
			tserv.voluntarilySacrificeAgent(agent)
			sacrificialElims[agentID] = agent
		}
		// ...plus 2*(n-v) random non-volunteers
		numNonVol := len(nonVolunteers)
		rand.Shuffle(numNonVol, func(i, j int) {
			nonVolunteers[i], nonVolunteers[j] = nonVolunteers[j], nonVolunteers[i]
		})

		for i := range min(numNonVol, 2*(neededVolunteers-actualVolunteers)) {
			agent := nonVolunteers[i]
			agentID := agent.GetID()
			tserv.involuntarilySacrificeAgent(agent)
			sacrificialElims[agentID] = agent
		}
	}

	return sacrificialElims
}

func updateAgentYsterofimia(agent infra.IExtendedAgent, agentsToRemove map[uuid.UUID]infra.IExtendedAgent, volunteerLookup map[uuid.UUID]struct{}) {
	networkEliminationCount := 0
	for friendID, esteem := range agent.GetNetwork() {
		// friend was not eliminated
		if _, ok := agentsToRemove[friendID]; !ok {
			continue
		}
		networkEliminationCount++
		ysterofimia := agent.GetYsterofimia()
		if _, ok := volunteerLookup[friendID]; ok {
			ysterofimia.IncrementSelfSacrificeCount()
			ysterofimia.AddSelfSacrificeEsteems(esteem)
		} else {
			ysterofimia.IncrementOtherEliminationCount()
			ysterofimia.AddOtherEliminationsEsteems(esteem)
		}

	}
	agent.IncrementNetworkEliminations(networkEliminationCount)
}

func (tserv *TMTServer) ApplyElimination() {
	tserv.LastEliminatedAgents = make([]infra.IExtendedAgent, 0)
	tserv.LastSelfSacrificedAgents = make([]infra.IExtendedAgent, 0)
	clusterEliminationCount := make(map[int]int) // number of eliminations per cluster
	agentsToRemove := make(map[uuid.UUID]infra.IExtendedAgent)

	tserv.updateAgentMortality()

	naturalElims := tserv.getNaturalEliminations()
	volunteers, nonVolunteers := tserv.stratifyVolunteers()
	sacrificialElims := tserv.getSacrificialEliminations(volunteers, nonVolunteers)

	// combine maps into one
	maps.Copy(agentsToRemove, naturalElims)
	maps.Copy(agentsToRemove, sacrificialElims)

	// also track eliminations per cluster and in network
	for agentID, agent := range agentsToRemove {
		clusterID := agent.GetClusterID()    // get the cluster ID of the agent
		clusterEliminationCount[clusterID]++ // increment the count for that cluster
		fmt.Print("Removing agent from server: ", agentID)
		tserv.RemoveAgent(agent)
	}

	// create hashset of volunteer IDs for ysterofimia
	volunteerLookup := make(map[uuid.UUID]struct{})
	for _, agent := range volunteers {
		agentID := agent.GetID()
		volunteerLookup[agentID] = struct{}{}
	}

	for _, agent := range tserv.GetAgentMap() {
		clusterID := agent.GetClusterID()
		if eliminatedInCluster, exists := clusterEliminationCount[clusterID]; exists {
			agent.IncrementClusterEliminations(eliminatedInCluster)
		}
		updateAgentYsterofimia(agent, agentsToRemove, volunteerLookup)
	}
}
