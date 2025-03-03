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
	NameID int

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

func CreateExtendedAgents(funcs agent.IExposedServerFunctions[infra.IExtendedAgent], configParam AgentConfig, grid *infra.Grid) *ExtendedAgent {
	return &ExtendedAgent{
		BaseAgent: agent.CreateBaseAgent(funcs),
		Server: funcs.(infra.IServer),
		NameID: 0,
		Attachment: []float32{rand.Float32(), rand.Float32()}, // Randomised anxiety and avoidance
		Network: make(map[uuid.UUID]float32), // Assign a unique UUID
		Age: rand.Intn(100), // Randomised age between 0 and 100
		MortalitySalience: false,
		SelfSacrificeWillingness: configParam.InitSacrificeWillingness,
		//ContextSacrifice: "",
		Position: [2]int{rand.Intn(grid.Width), rand.Intn(grid.Height)},
	}
}


func (ea *ExtendedAgent) GetName() int {
    return ea.NameID
}

func (ea *ExtendedAgent) SetName(name int) {
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
	return ea.Network
}

func (ea *ExtendedAgent) SetNetwork(network map[uuid.UUID]float32) {
	ea.Network = network
}

// distance between two agents on grid
func (ea *ExtendedAgent) DistanceTo(other *ExtendedAgent) float64 {
	return infra.Distance(ea.Position, other.Position)
}

func (a *ExtendedAgent) AddRelationship(otherID uuid.UUID, strength float32) {
	if _, exists := a.Network[otherID]; !exists {
		a.Network[otherID] = strength
	}
}

func (a *ExtendedAgent) UpdateRelationship(otherID uuid.UUID, change float32) {
	if _, exists := a.Network[otherID]; exists {
		a.Network[otherID] += change

		// Keep values between 0 and 1
		if a.Network[otherID] > 1 {
			a.Network[otherID] = 1
		} else if a.Network[otherID] < 0 {
			a.Network[otherID] = 0
		}
	}
}

func (a *ExtendedAgent) DecayRelationships() {
	for id := range a.Network {
		a.Network[id] -= 0.05 // Reduce strength over time
		if a.Network[id] < 0 {
			delete(a.Network, id) // Remove weak relationships
		}
	}
}

func (a *ExtendedAgent) Interact(other *ExtendedAgent) {
	a.UpdateRelationship(other.GetID(), 0.1)
	other.UpdateRelationship(a.GetID(), 0.1)
}

func (ea *ExtendedAgent) GetAge() int{
    return ea.Age
}

func (ea *ExtendedAgent) SetAge(age int) {
    ea.Age = age
}

func (a *ExtendedAgent) MoveRandomly(grid *infra.Grid) {
	a.Position = [2]int{rand.Intn(grid.Width), rand.Intn(grid.Height)}
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


// func (ea *ExtendedAgent) GetContextSacrifice() string {
//     return ea.ContextSacrifice
// }

// func (ea *ExtendedAgent) SetContextSacrifice(context string) {
//     ea.ContextSacrifice = context
// }


// Decision-making logic
func (ea *ExtendedAgent) DecideSacrifice() float32 {
    //TO-DO: Fuzzy logic stuff

	
    //fmt.Printf("Agent %d decision: %v\n", a.NameID, a.SacrificeChoice)
    
	fmt.Printf("Agent %d decided to %s \n",
        ea.NameID,
        map[float32]string{}[ea.SelfSacrificeWillingness])
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

