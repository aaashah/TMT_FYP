package server

import (
	"math/rand"

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

func (tserv *TMTServer) voluntarilySacrificeAgent(agent infra.IExtendedAgent) {
	pos := agent.GetPosition()
	tserv.Grid.PlaceTemple(pos.X, pos.Y)

	tserv.lastEliminatedAgents = append(tserv.lastEliminatedAgents, agent)
	tserv.lastSelfSacrificedAgents = append(tserv.lastSelfSacrificedAgents, agent)
	// fmt.Printf("Agent %v has been eliminated (voluntary)\n", agent.GetID())
}

func (tserv *TMTServer) involuntarilySacrificeAgent(agent infra.IExtendedAgent) {
	pos := agent.GetPosition()
	tserv.Grid.PlaceTombstone(pos.X, pos.Y)
	tserv.lastEliminatedAgents = append(tserv.lastEliminatedAgents, agent)
	// fmt.Printf("Agent %v has been eliminated (non-voluntary)\n", agent.GetID())
}

func (tserv *TMTServer) getNaturalEliminationReport() map[uuid.UUID]infra.DeathInfo {
	naturalReport := make(map[uuid.UUID]infra.DeathInfo)
	for agentID, agent := range tserv.GetAgentMap() {
		if !agent.IsAlive() {
			naturalReport[agentID] = infra.DeathInfo{Agent: agent, WasVoluntary: false}
		}
	}
	return naturalReport
}

func (tserv *TMTServer) stratifyVolunteers() ([]infra.IExtendedAgent, []infra.IExtendedAgent) {
	volunteers := make([]infra.IExtendedAgent, 0)
	nonVolunteers := make([]infra.IExtendedAgent, 0)

	// Separate volunteers and non-volunteers
	for _, agent := range tserv.GetAgentMap() {
		// don't allow naturally-dead agents to sacrifice
		if !agent.IsAlive() {
			continue
		}
		// check if agent volunteered
		if agent.GetASPDecision(tserv.Grid) == infra.SELF_SACRIFICE {
			volunteers = append(volunteers, agent)
		} else {
			nonVolunteers = append(nonVolunteers, agent)
		}
	}
	return volunteers, nonVolunteers
}

func (tserv *TMTServer) getSacrificialEliminationReport() map[uuid.UUID]infra.DeathInfo {
	volunteers, nonVolunteers := tserv.stratifyVolunteers()
	sacrificialReport := make(map[uuid.UUID]infra.DeathInfo)

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
			sacrificialReport[agentID] = infra.DeathInfo{Agent: agent, WasVoluntary: true}
		}
	} else {
		//eliminate all volunteers...
		for _, agent := range volunteers {
			agentID := agent.GetID()
			sacrificialReport[agentID] = infra.DeathInfo{Agent: agent, WasVoluntary: true}
		}
		// ...plus 2*(n-v) random non-volunteers
		numNonVol := len(nonVolunteers)
		rand.Shuffle(numNonVol, func(i, j int) {
			nonVolunteers[i], nonVolunteers[j] = nonVolunteers[j], nonVolunteers[i]
		})

		for i := range min(numNonVol, 2*(neededVolunteers-actualVolunteers)) {
			agent := nonVolunteers[i]
			agentID := agent.GetID()
			sacrificialReport[agentID] = infra.DeathInfo{Agent: agent, WasVoluntary: false}
		}
	}

	return sacrificialReport
}

func (tserv *TMTServer) updateAgentYsterofimia(deathReport map[uuid.UUID]infra.DeathInfo) {
	for _, agent := range tserv.GetAgentMap() {
		networkEliminationCount := 0
		for friendID, esteem := range agent.GetNetwork() {
			// friend was eliminated (found in death report)
			if deathInfo, dead := deathReport[friendID]; dead {
				networkEliminationCount++
				ysterofimia := agent.GetYsterofimia()
				if deathInfo.WasVoluntary {
					ysterofimia.IncrementSelfSacrificeCount()
					ysterofimia.AddSelfSacrificeEsteems(esteem)
				} else {
					ysterofimia.IncrementOtherEliminationCount()
					ysterofimia.AddOtherEliminationsEsteems(esteem)
				}
			}
		}
		agent.IncrementNetworkEliminations(networkEliminationCount)
	}
}

func (tserv *TMTServer) updateClusterEliminations(deathReport map[uuid.UUID]infra.DeathInfo) {
	counts := make(map[int]int)
	for _, deathInfo := range deathReport {
		clusterID := deathInfo.Agent.GetClusterID()
		counts[clusterID]++
	}
	for _, agent := range tserv.GetAgentMap() {
		clusterID := agent.GetClusterID()
		if count, exists := counts[clusterID]; exists {
			agent.IncrementClusterEliminations(count)
		}
	}
}

func (tserv *TMTServer) updateAgentHeroism(deathReport map[uuid.UUID]infra.DeathInfo) {
	for _, deathInfo := range deathReport {
		agent := deathInfo.Agent
		if deathInfo.WasVoluntary {
			agent.IncrementHeroism()
		}
	}
}

func (tserv *TMTServer) applyElimination(deathReport map[uuid.UUID]infra.DeathInfo) {
	for _, deathInfo := range deathReport {
		deadAgent := deathInfo.Agent
		tserv.RemoveAgent(deadAgent)
	}
}

func (tserv *TMTServer) performSacrifices(deathReport map[uuid.UUID]infra.DeathInfo) {
	tserv.lastEliminatedAgents = nil
	tserv.lastSelfSacrificedAgents = nil

	for _, deathInfo := range deathReport {
		deadAgent := deathInfo.Agent
		if deathInfo.WasVoluntary {
			tserv.voluntarilySacrificeAgent(deadAgent)
		} else {
			tserv.involuntarilySacrificeAgent(deadAgent)
		}
	}
}
