package main

import (
	"github.com/MattSScott/TMT_SOMAS/agents"
	"github.com/MattSScott/TMT_SOMAS/config"
	"github.com/MattSScott/TMT_SOMAS/infra"
	"github.com/MattSScott/TMT_SOMAS/server"
	"github.com/google/uuid"
)

func main() {
	config := config.NewConfig()
	serv := server.CreateTMTServer(config)
	serv.SetGameRunner(serv)

	parent1ID, parent2ID := uuid.Nil, uuid.Nil
	agentPopulation := make([]infra.IExtendedAgent, 0)

	for i := 0; i < config.NumAgents; i += 4 {
		agentPopulation = append(agentPopulation, agents.CreateSecureAgent(serv, parent1ID, parent2ID))
		agentPopulation = append(agentPopulation, agents.CreateDismissiveAgent(serv, parent1ID, parent2ID))
		agentPopulation = append(agentPopulation, agents.CreatePreoccupiedAgent(serv, parent1ID, parent2ID))
		agentPopulation = append(agentPopulation, agents.CreateFearfulAgent(serv, parent1ID, parent2ID))
	}

	// Set probability p for Erdős–Rényi network
	for _, agent := range agentPopulation {
		serv.AddAgent(agent)
		if config.Debug {
			agent.AgentInitialised() // Call the method to print agent details
		}
	}

	// Initialize social network after agents are created
	serv.InitialiseRandomNetwork(config.ConnectionProbability)

	// Start server
	serv.Start()
}
