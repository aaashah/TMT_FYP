package tests

import (
	"math/rand"
	"testing"

	//"github.com/google/uuid"
	"github.com/aaashah/TMT_FYP/agents"
	"github.com/aaashah/TMT_FYP/config"
	"github.com/aaashah/TMT_FYP/infra"
	"github.com/aaashah/TMT_FYP/server"
)

// TestFindClosestFriend verifies the logic of finding the closest friend in an agent's network.
func TestFindClosestFriend(t *testing.T) {
	// Initialize server and config
	conf := config.NewConfig()
	serv := server.CreateTMTServer(conf)

	// Create the main agent
	agent := agents.CreateSecureAgent(serv)

	// Create friends
	friend1 := agents.CreateSecureAgent(serv)
	friend2 := agents.CreateSecureAgent(serv)
	friend3 := agents.CreateSecureAgent(serv)
	serv.AddAgent(agent)
	serv.AddAgent(friend1)
	serv.AddAgent(friend2)
	serv.AddAgent(friend3)

	// Set positions
	agent.SetPosition(infra.PositionVector{X: 0, Y: 0})
	friend1.SetPosition(infra.PositionVector{X: 1, Y: 1}) // distance ~1.41
	friend2.SetPosition(infra.PositionVector{X: 2, Y: 2}) // distance ~2.82
	friend3.SetPosition(infra.PositionVector{X: -1, Y: -1}) // distance ~1.41

	// Add friends to the agent's social network (weights arbitrary)
	serv.CreateNetworkConnection(agent.GetID(), agent.GetID(), 1.0)
	serv.CreateNetworkConnection(agent.GetID(), friend1.GetID(), 1.0)
	serv.CreateNetworkConnection(agent.GetID(), friend2.GetID(), 1.0)
	serv.CreateNetworkConnection(agent.GetID(), friend3.GetID(), 1.0)
	agent.AddToSocialNetwork(agent.GetID(), 1.0)
	agent.AddToSocialNetwork(friend1.GetID(), 1.0)
	agent.AddToSocialNetwork(friend2.GetID(), 1.0)
	agent.AddToSocialNetwork(friend3.GetID(), 1.0)

	// Seed random for deterministic selection
	rand.Seed(42)

	// Test: should return either friend1 or friend3 (both same distance)
	closest := agent.FindClosestFriend()
	if closest == nil {
		t.Fatalf("Expected closest friend, got nil")
	}
	closestID := closest.GetID()
	if closestID != friend1.GetID() && closestID != friend3.GetID() {
		t.Errorf("Expected closest to be friend1 or friend3, got: %v", closestID)
	}

	// Remove friend3 from network, now only friend1 is closest
	serv.SeverNetworkConnection(agent.GetID(), friend3.GetID())

	closest = agent.FindClosestFriend()
	if closest.GetID() != friend1.GetID() {
		t.Errorf("Expected closest friend1, got: %v", closest.GetID())
	}

	// Remove all friends: should return nil
	serv.SeverNetworkConnection(agent.GetID(), friend1.GetID())
	serv.SeverNetworkConnection(agent.GetID(), friend2.GetID())
	closest = agent.FindClosestFriend()
	if closest != nil {
		t.Errorf("Expected no closest friend (nil), got: %v", closest.GetID())
	}

	// Add only one friend, should return that friend
	agent.GetNetwork()[friend2.GetID()] = 1.0
	closest = agent.FindClosestFriend()
	if closest.GetID() != friend2.GetID() {
		t.Errorf("Expected closest to be friend2, got: %v", closest.GetID())
	}
}