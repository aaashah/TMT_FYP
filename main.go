package main

import (
	"github.com/aaashah/TMT_FYP/agents"
	"github.com/aaashah/TMT_FYP/config"
	"github.com/aaashah/TMT_FYP/infra"
	"github.com/aaashah/TMT_FYP/server"
)

func main() {
	config := config.NewConfig()
	serv := server.CreateTMTServer(config)
	serv.SetGameRunner(serv)

	agentPopulation := make([]infra.IExtendedAgent, 0)

	totalAgents := float64(config.NumAgents)
	for range int(totalAgents * config.DismissiveProp) {
		agentPopulation = append(agentPopulation, agents.CreateDismissiveAgent(serv))
	}
	for range int(totalAgents * config.FearfulProp) {
		agentPopulation = append(agentPopulation, agents.CreateFearfulAgent(serv))
	}
	for range int(totalAgents * config.PreoccupiedProp) {
		agentPopulation = append(agentPopulation, agents.CreatePreoccupiedAgent(serv))
	}
	for range int(totalAgents * config.SecureProp) {
		agentPopulation = append(agentPopulation, agents.CreateSecureAgent(serv))
	}

	for _, agent := range agentPopulation {
		serv.AddAgent(agent)
		if config.Debug {
			agent.AgentInitialised()
		}
	}

	serv.Start()
}
