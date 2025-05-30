package tests

import (
	"testing"

	"github.com/aaashah/TMT_FYP/agents"
	"github.com/aaashah/TMT_FYP/config"

	"github.com/aaashah/TMT_FYP/infra"
	"github.com/aaashah/TMT_FYP/server"
	"github.com/stretchr/testify/assert"
	//"github.com/google/uuid"
)

// Testing relative age to network

func TestRelativeAgeToNetwork(t *testing.T) {
	// Initialize server and config
	conf := config.NewConfig()
	serv := server.CreateTMTServer(conf)

	// Create the main agent
	agent := agents.CreateSecureAgent(serv)
	serv.AddAgent(agent)

	// Create friends with different ages
	friend1 := agents.CreateSecureAgent(serv)
	friend2 := agents.CreateSecureAgent(serv)
	friend3 := agents.CreateSecureAgent(serv)

	//new telomeres for friends
	friend1Age := 20 // 20 years old
	for agent.GetAge() < friend1Age {
		agent.IncrementAge()
	}

	friend2Age := 30 // 30 years old
	for agent.GetAge() < friend2Age {
		agent.IncrementAge()
	}

	friend3Age := 40 // 40 years old
	for agent.GetAge() < friend3Age {
		agent.IncrementAge()
	}

	// Add friends to the agent's social network (weights arbitrary)
	serv.CreateNetworkConnection(agent.GetID(), friend1.GetID(), 1.0)
	serv.CreateNetworkConnection(agent.GetID(), friend2.GetID(), 1.0)
	serv.CreateNetworkConnection(agent.GetID(), friend3.GetID(), 1.0)

	agent.AddToSocialNetwork(friend1.GetID(), 1.0)
	agent.AddToSocialNetwork(friend2.GetID(), 1.0)
	agent.AddToSocialNetwork(friend3.GetID(), 1.0)

	// Test: should return the average age of friends in the network
	expectedAvgAge := (20 + 30 + 40) / 3.0
	assert.Equal(t, expectedAvgAge, agent.RelativeAgeToNetwork())
}

// Testing memorial proximity
func TestMemorialProximity(t *testing.T) {
	// Setup
	conf := config.NewConfig()
	serv := server.CreateTMTServer(conf)

	// Create the main agent
	agent := agents.CreateSecureAgent(serv)
	serv.AddAgent(agent)
	agent.SetPosition(infra.PositionVector{X: 0, Y: 0})
	agent.SetClusterID(1)

	// Create grid with one tombstone at (1, 0)
	grid := &infra.Grid{
		Tombstones: []infra.PositionVector{{X: 1, Y: 0}},
		Temples:    []infra.PositionVector{},
	}

	// No other agents in the cluster
	proximity := agent.GetMemorialProximity(grid)
	assert.Equal(t, float32(0.0), proximity, "No cluster members, only memorials → proximity 0")

	// Add another agent in the same cluster
	friend := agents.CreateSecureAgent(serv)
	serv.AddAgent(friend)
	friend.SetPosition(infra.PositionVector{X: 0, Y: 1})
	friend.SetClusterID(1)

	proximity = agent.GetMemorialProximity(grid)
	// compute expected influences
	// distance to tombstone = 1
	memorialInfluence := 1.0

	// distance to friend = 1
	clusterInfluence := 1.0

	expectedProximity := float32(clusterInfluence) / float32(clusterInfluence+memorialInfluence)
	assert.InDelta(t, expectedProximity, proximity, 1e-6, "Equal influence from friend and memorial")

	// Remove memorial, test only cluster influence
	grid.Tombstones = []infra.PositionVector{}
	proximity = agent.GetMemorialProximity(grid)
	assert.Equal(t, float32(1.0), proximity, "Only cluster influence → proximity 1")

	// Remove friend from cluster (different clusterID)
	friend.SetClusterID(2)
	proximity = agent.GetMemorialProximity(grid)
	assert.Equal(t, float32(0.0), proximity, "No cluster members or memorials → proximity 0")
}