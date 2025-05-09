package main

import (
	"github.com/MattSScott/TMT_SOMAS/agents"
	"github.com/MattSScott/TMT_SOMAS/config"
	"github.com/MattSScott/TMT_SOMAS/infra"
	"github.com/MattSScott/TMT_SOMAS/server"
)

func main() {
	config := config.NewConfig()
	serv := server.CreateTMTServer(config)
	serv.SetGameRunner(serv)

	agentPopulation := make([]infra.IExtendedAgent, 0)

	for i := 0; i < config.NumAgents; i += 4 {
		agentPopulation = append(agentPopulation, agents.CreateSecureAgent(serv))
		agentPopulation = append(agentPopulation, agents.CreateDismissiveAgent(serv))
		agentPopulation = append(agentPopulation, agents.CreatePreoccupiedAgent(serv))
		agentPopulation = append(agentPopulation, agents.CreateFearfulAgent(serv))
	}

	for _, agent := range agentPopulation {
		serv.AddAgent(agent)
		if config.Debug {
			agent.AgentInitialised()
		}
	}

	// Start server
	serv.Start()
}
