package tests

import (
	"testing"

	"github.com/aaashah/TMT_FYP/agents"
	"github.com/aaashah/TMT_FYP/config"
	"github.com/aaashah/TMT_FYP/infra"
	"github.com/aaashah/TMT_FYP/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInitialiseRandomNetwork(t *testing.T) {
	// Initialize server and config
	conf := config.NewConfig()
	serv := server.CreateTMTServer(conf)

	// Create the main agent and set position
	agent := agents.CreateSecureAgent(serv)
	agent.SetPosition(infra.PositionVector{X: 0, Y: 0})

	// create some friends
	for i := 0; i < 5; i++ {
		friend := agents.CreateSecureAgent(serv)
		friend.SetPosition(infra.PositionVector{X: int(i), Y: int(i)})
	}
	
	serv.InitialiseRandomNetworkForAgent(agent)
	network := agent.GetNetwork()

	// Check that the agent has at least one connection
	assert.Greater(t, len(network), 0, "Agent should have at least one friend")

	// Check that the agent is not connected to itself
	_, selfExists := network[agent.GetID()]
	assert.True(t, selfExists, "Agent should be connected to itself")

	// Check that all entries are unique
	seen := make(map[uuid.UUID]bool)
	for id := range network {
		if seen[id] {
			t.Errorf("Duplicate friend ID found: %v", id)
		}
		seen[id] = true
	}
}