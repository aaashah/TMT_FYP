package main

import (
	"math/rand"

	"github.com/aaashah/TMT_Attachment/agents"
	"github.com/aaashah/TMT_Attachment/config"
	"github.com/aaashah/TMT_Attachment/infra"
	"github.com/aaashah/TMT_Attachment/server"
	"github.com/google/uuid"
)

func main() {
	config := config.NewConfig()
	serv := server.CreateTMTServer(config)
	serv.SetGameRunner(serv)

	parent1ID, parent2ID := uuid.Nil, uuid.Nil
	agentPopulation := make([]infra.IExtendedAgent, 0)

	for i := 0; i < config.NumAgents; i += 4 {
		agentPopulation = append(agentPopulation, agents.CreateSecureAgent(serv, parent1ID, parent2ID, rand.Uint32()))
		agentPopulation = append(agentPopulation, agents.CreateDismissiveAgent(serv, parent1ID, parent2ID, rand.Uint32()))
		agentPopulation = append(agentPopulation, agents.CreatePreoccupiedAgent(serv, parent1ID, parent2ID, rand.Uint32()))
		agentPopulation = append(agentPopulation, agents.CreateFearfulAgent(serv, parent1ID, parent2ID, rand.Uint32()))
	}

	// Set probability p for Erdős–Rényi network
	for _, agent := range agentPopulation {
		serv.AddAgent(agent)
		if config.Debug {
			agent.AgentInitialised() // Call the method to print agent details
		}
		//fmt.Printf("Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", agent.GetID(), agent.GetAge(), agent.GetAttachment()[0], agent.GetAttachment()[1])
	}

	// Initialize social network after agents are created
	serv.InitialiseRandomNetwork(config.ConnectionProbability)

	// Start server
	serv.Start()
}
