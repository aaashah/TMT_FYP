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

	Age int
	AgeA int // Age where mortality probability starts increasing
	AgeB int // Age where agent is definitely eliminated
	Telomere float32 // Determines lifespan decay (death probability)

	Position [2]int
	MovementPolicy string // Defines how movement is determined

	//History Tracking
	ObservedEliminationsCluster int
	ObservedEliminationsNetwork int
	Heroism                     int // Number of voluntary self-sacrifices

	// Social network and kinship group
	Network map[uuid.UUID]float32 // stores relationship strengths
	KinshipGroup        []uuid.UUID  // Descendants 

	Attachment []float32 // Attachment orientations: [anxiety, avoidance].

	// **Decision-Making Parameters**
	ASP map[string]float64 // Parameters for decision-making
	PTS map[string]float64 // Parameters for behavior protocols

	// **Worldview (32-bit opinion vector)**
	Worldview [32]bool

	// **Isterofimia (Posthumous Recognition)**
	Isterofimia float64 // Observation of self-sacrifice vs self-preservation

	Mortality bool

	MortalitySalience float32 //section in ASP module
	WorldviewValidation float32 //section in ASP module
	RelationshipValidation float32 //section in ASP module

	SelfSacrificeWillingness float32 //ASP result
}


type AgentConfig struct {
	InitSacrificeWillingness float32
}
//var _ infra.IExtendedAgent = (*ExtendedAgent)(nil)

func CreateExtendedAgents(server infra.IServer, configParam AgentConfig, grid *infra.Grid) *ExtendedAgent {
	A := rand.Intn(25) + 40  // (40-65)
	B := A + rand.Intn(35) + 20  // Random max age (60 - 100)

	return &ExtendedAgent{
		BaseAgent: agent.CreateBaseAgent(server),
		Server:    server,
		NameID:    uuid.New(),
		Attachment: []float32{rand.Float32(), rand.Float32()}, // Randomized anxiety and avoidance
		Network:    make(map[uuid.UUID]float32),
		Age:        rand.Intn(50),
		AgeA:       A,
		AgeB:       B,
		Mortality: false,
		SelfSacrificeWillingness: configParam.InitSacrificeWillingness,
		Position: [2]int{rand.Intn(grid.Width) + 1, rand.Intn(grid.Height) + 1},
	}
}


func (ea *ExtendedAgent) GetName() uuid.UUID {
	return ea.GetID()
}

func (ea *ExtendedAgent) SetName(name uuid.UUID) {
    ea.NameID = ea.GetID()
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

func (ea *ExtendedAgent) Move(grid *infra.Grid) {
	newX, newY := grid.GetValidMove(ea.Position[0], ea.Position[1]) // Get a valid move
	grid.UpdateAgentPosition(ea, newX, newY)    // Update position in the grid
	ea.Position = [2]int{newX, newY}             // ✅ Assign new position
	fmt.Printf("Agent %v moved to (%d, %d)\n", ea.GetID(), newX, newY)
}

// distance between two agents on grid
// func (ea *ExtendedAgent) DistanceTo(other *ExtendedAgent) float64 {
// 	return infra.Distance(ea.Position, other.Position)
// }

func (ea *ExtendedAgent) AddRelationship(otherID uuid.UUID, strength float32) {
	ea.Server.UpdateAgentRelationship(ea.GetID(), otherID, strength)
}

func (ea *ExtendedAgent) UpdateRelationship(otherID uuid.UUID, change float32) {
	ea.Server.UpdateAgentRelationship(ea.GetID(), otherID, change)
}


// GetAge generates an age following a beta-like distribution approximating the UK population.

func (ea *ExtendedAgent) GetAge() int {
	// Beta distribution parameters (adjusted to fit UK population shape)
	return ea.Age
}

func (ea *ExtendedAgent) GetTelomere() float32 {

	if ea.Age < ea.AgeA {
		return 0.005 * float32(ea.Age) // Small increasing probability
	} else if ea.Age >= ea.AgeB {
		return 1.0 // Guaranteed death at AgeB
	} else {
		// Linearly increasing probability from AgeA to AgeB
		return float32(ea.Age-ea.AgeA) / float32(ea.AgeB-ea.AgeA)
	}
}


func (ea *ExtendedAgent) SetAge(age int) {
    ea.Age = age
}

// GetMortality returns the mortality status of the agent.
func (ea *ExtendedAgent) GetMortality() bool {
	probDeath := ea.GetTelomere() // get probability of death
	randVal := rand.Float32()     // random value between 0 and 1
	ea.Mortality = randVal < probDeath // Random chance to die of natural causes
	return ea.Mortality
}

func (ea *ExtendedAgent) IsMortalitySalient() bool {
    return ea.Mortality
}

func (ea *ExtendedAgent) SetMortalitySalience(ms bool) {
    ea.Mortality = ms
}

func (ea *ExtendedAgent) GetSelfSacrificeWillingness() float32 {
	return ea.SelfSacrificeWillingness
}

// Decision-making logic
func (ea *ExtendedAgent) DecideSacrifice() float32 {
    //TO-DO: Fuzzy logic stuff

	ea.SelfSacrificeWillingness = rand.Float32() // Random willingness to sacrifice

	
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
	//fmt.Printf("[DEBUG] Fetching Age in RecordAgentStatus: %d for Agent %v\n", instance.GetAge(), instance.GetID()) 
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
