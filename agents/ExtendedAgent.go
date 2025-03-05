package agents

import (
	"fmt"
	"math/rand"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)


type ExtendedAgent struct {
	*agent.BaseAgent[infra.IExtendedAgent]
	Server infra.IServer
	NameID uuid.UUID

	//private
	Network map[uuid.UUID]float32 // stores relationship strengths
	Attachment []float32 // Attachment orientations: [anxiety, avoidance].
	Age int
	MortalitySalience bool
	//WorldviewValidation float32
	//RelationshipValidation float32


	// dynamic
	SelfSacrificeWillingness float32
	//ContextSacrifice string
	Position [2]int
}


type AgentConfig struct {
	InitSacrificeWillingness float32
}
//var _ infra.IExtendedAgent = (*ExtendedAgent)(nil)

func CreateExtendedAgents(server infra.IServer, configParam AgentConfig, grid *infra.Grid) *ExtendedAgent {
	return &ExtendedAgent{
		BaseAgent: agent.CreateBaseAgent(server),
		Server:    server,
		NameID:    uuid.New(),
		Attachment: []float32{rand.Float32(), rand.Float32()}, // Randomized anxiety and avoidance
		Network:    make(map[uuid.UUID]float32),
		Age:        rand.Intn(100),
		MortalitySalience: false,
		SelfSacrificeWillingness: configParam.InitSacrificeWillingness,
		Position: [2]int{rand.Intn(grid.Width), rand.Intn(grid.Height)},
	}
}


func (ea *ExtendedAgent) GetName() uuid.UUID {
    return ea.NameID
}

func (ea *ExtendedAgent) SetName(name uuid.UUID) {
    ea.NameID = name
}

func (ea *ExtendedAgent) GetPosition() [2]int {
	return ea.Position
}

func (ea *ExtendedAgent) GetAttachment() []float32 {
    return ea.Attachment
}

func (ea *ExtendedAgent) SetAttachment(attachment []float32) {
    if len(attachment) != 2 {
        panic("Attachment must have exactly two elements: [anxiety, avoidance]")
    }
    ea.Attachment = attachment
}


func (ea *ExtendedAgent) GetNetwork() map[uuid.UUID]float32 {
	updatedNetwork := make(map[uuid.UUID]float32)
	for id, strength := range ea.Network {
		updatedNetwork[id] = strength
	}
	return updatedNetwork
}

func (ea *ExtendedAgent) SetNetwork(network map[uuid.UUID]float32) {
	ea.Network = network
}

// distance between two agents on grid
func (ea *ExtendedAgent) DistanceTo(other *ExtendedAgent) float64 {
	return infra.Distance(ea.Position, other.Position)
}

func (ea *ExtendedAgent) AddRelationship(otherID uuid.UUID, strength float32) {
	ea.Server.UpdateAgentRelationship(ea.NameID, otherID, strength)
}

func (ea *ExtendedAgent) UpdateRelationship(otherID uuid.UUID, change float32) {
	ea.Server.UpdateAgentRelationship(ea.NameID, otherID, change)
}


// Moves an agent towards the strongest connection in its network.
func (ea *ExtendedAgent) MoveAttractedToNetwork(grid *infra.Grid, server infra.IServer) {
	if len(ea.Network) == 0 {
		// No social ties → move randomly
		ea.Position[0] += rand.Intn(3) - 1
		ea.Position[1] += rand.Intn(3) - 1
		return
	}

	// Find the most attractive agent(s)
	var bestNeighbor uuid.UUID
	maxAttraction := float32(-1)

	for neighborID, strength := range ea.Network {
		if strength > maxAttraction {
			bestNeighbor = neighborID
			maxAttraction = strength
		}
	}

	if bestNeighbor == uuid.Nil {
		return // No valid movement target
	}

	// Get bestNeighbor's position from server
	bestPos, exists := server.GetAgentPosition(bestNeighbor)
	if !exists {
		return
	}

	// Compute movement direction
	dx := bestPos[0] - ea.Position[0]
	dy := bestPos[1] - ea.Position[1]

	moveX, moveY := 0, 0
	if dx > 0 {
		moveX = 1
	} else if dx < 0 {
		moveX = -1
	}

	if dy > 0 {
		moveY = 1
	} else if dy < 0 {
		moveY = -1
	}

	// Move 1 or 2 steps in direction
	step := rand.Intn(2) + 1
	newX := ea.Position[0] + moveX*step
	newY := ea.Position[1] + moveY*step

	// Keep inside grid bounds
	if newX < 0 {
		newX = 0
	} else if newX >= grid.Width {
		newX = grid.Width - 1
	}

	if newY < 0 {
		newY = 0
	} else if newY >= grid.Height {
		newY = grid.Height - 1
	}

	// Update position
	ea.Position = [2]int{newX, newY}
	fmt.Printf("Agent %v moved to (%d, %d) towards %v\n", ea.NameID, newX, newY, bestNeighbor)
}

func (ea *ExtendedAgent) GetAge() int{
    return ea.Age
}

func (ea *ExtendedAgent) SetAge(age int) {
    ea.Age = age
}

func (ea *ExtendedAgent) IsMortalitySalient() bool {
    return ea.MortalitySalience
}

func (ea *ExtendedAgent) SetMortalitySalience(ms bool) {
    ea.MortalitySalience = ms
}

func (ea *ExtendedAgent) GetSelfSacrificeWillingness() float32 {
	return ea.SelfSacrificeWillingness
}

// Decision-making logic
func (ea *ExtendedAgent) DecideSacrifice() float32 {
    //TO-DO: Fuzzy logic stuff

	ea.SelfSacrificeWillingness = rand.Float32()

	
    //fmt.Printf("Agent %d decision: %v\n", a.NameID, a.SacrificeChoice)
    
	// fmt.Printf("Agent %v willing to sacrifice by %s \n",
    //     ea.NameID,
    //     map[float32]string{}[ea.SelfSacrificeWillingness])
    return ea.SelfSacrificeWillingness
}

func (ea *ExtendedAgent) GetExposedInfo() infra.ExposedAgentInfo {
	return infra.ExposedAgentInfo{
		AgentUUID: ea.GetID(),
	}
}

// ----------------------- Data Recording Functions -----------------------
func (mi *ExtendedAgent) RecordAgentStatus(instance infra.IExtendedAgent) gameRecorder.AgentRecord {
	record := gameRecorder.NewAgentRecord(
		instance.GetID(),
		instance.GetAge(),
		instance.GetPosition()[0],
		instance.GetPosition()[1],
		instance.GetSelfSacrificeWillingness(),
		instance.GetAttachment()[0],
		instance.GetAttachment()[1],
		"1",
		//instance.GetWorldviewValidation(),
		//instance.GetRelationshipValidation(),
		
	)
	return record
}
