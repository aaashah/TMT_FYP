package main

import (
	"math/rand"

	agents "github.com/aaashah/TMT_Attachment/agents"
	infra "github.com/aaashah/TMT_Attachment/infra"
	tmtServer "github.com/aaashah/TMT_Attachment/server"
	"github.com/google/uuid"
)

func main() {
	serv := tmtServer.CreateTMTServer()
	serv.SetGameRunner(serv)

	const numAgents int = 40
	parent1ID, parent2ID := uuid.Nil, uuid.Nil
	agentPopulation := make([]infra.IExtendedAgent, numAgents)

	for i := 0; i < numAgents; i += 4 {
		agentPopulation = append(agentPopulation, agents.CreateSecureAgent(serv, parent1ID, parent2ID, rand.Uint32()))
		agentPopulation = append(agentPopulation, agents.CreateDismissiveAgent(serv, parent1ID, parent2ID, rand.Uint32()))
		agentPopulation = append(agentPopulation, agents.CreatePreoccupiedAgent(serv, parent1ID, parent2ID, rand.Uint32()))
		agentPopulation = append(agentPopulation, agents.CreateFearfulAgent(serv, parent1ID, parent2ID, rand.Uint32()))
	}

	// Set probability p for Erdős–Rényi network
	for _, agent := range agentPopulation {
		serv.AddAgent(agent)
		agent.AgentInitialised() // Call the method to print agent details
		//fmt.Printf("Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", agent.GetID(), agent.GetAge(), agent.GetAttachment()[0], agent.GetAttachment()[1])
	}

	// Initialize social network after agents are created
	const connectionProbability = 0.35
	serv.InitialiseRandomNetwork(connectionProbability)

	// Start server
	serv.Start()
}
