package server

import (
	"math/rand"
	"testing"

	"github.com/aaashah/TMT_FYP/infra"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRunKMeansBasicTwoClusters(t *testing.T) {
	// Seed random for deterministic centroid selection
	rand.Seed(42)

	// Two points far apart → should be in different clusters if numClusters=2
	positionMap := make(map[uuid.UUID]infra.PositionVector)
	agent1ID := uuid.New()
	agent2ID := uuid.New()

	positionMap[agent1ID] = infra.PositionVector{X: 0, Y: 0}
	positionMap[agent2ID] = infra.PositionVector{X: 100, Y: 100}

	assignments := runKMeans(positionMap, 2)

	assert.Equal(t, 2, len(assignments), "Both agents should be assigned to a cluster")
	assert.NotEqual(t, assignments[agent1ID], assignments[agent2ID], "Agents far apart should be in different clusters")
}

func TestRunKMeansSingleCluster(t *testing.T) {
	rand.Seed(42)
	positionMap := make(map[uuid.UUID]infra.PositionVector)

	for i := 0; i < 5; i++ {
		id := uuid.New()
		positionMap[id] = infra.PositionVector{X: int(i), Y: int(i)}
	}

	assignments := runKMeans(positionMap, 1)

	for _, cluster := range assignments {
		assert.Equal(t, 0, cluster, "All agents should be in cluster 0")
	}
}

func TestRunKMeansMoreClustersThanPoints(t *testing.T) {
	rand.Seed(42)
	positionMap := make(map[uuid.UUID]infra.PositionVector)

	agentID := uuid.New()
	positionMap[agentID] = infra.PositionVector{X: 10, Y: 10}

	assignments := runKMeans(positionMap, 3)

	assert.Equal(t, 1, len(assignments), "Only one agent to assign")
}

func TestRunKMeansEmptyInput(t *testing.T) {
	positionMap := make(map[uuid.UUID]infra.PositionVector)

	assignments := runKMeans(positionMap, 2)
	assert.Nil(t, assignments, "No input positions → should return nil")
}