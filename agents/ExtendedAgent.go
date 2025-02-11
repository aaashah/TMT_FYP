package agents

import (
	"fmt"
	"math/rand"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)


type ExtendedAgent struct {
	*agent.BaseAgent[infra.IExtendedAgent]
	Server infra.IServer
	NameID int

	//private
	Network map[uuid.UUID]int
	Attachment []float32 // Attachment orientations: [anxiety, avoidance].
	Age int
	MortalitySalience bool

	// dynamic
	SacrificeChoice bool
	ContextSacrifice string
}


type AgentConfig struct {
	InitSacrificeChoice bool
}

func CreateExtendedAgents(funcs agent.IExposedServerFunctions[infra.IExtendedAgent], configParam AgentConfig) *ExtendedAgent {
	return &ExtendedAgent{
		BaseAgent: agent.CreateBaseAgent(funcs),
		Server: funcs.(infra.IServer),
		NameID: 0,
		Attachment: []float32{rand.Float32(), rand.Float32()}, // Randomised anxiety and avoidance
		Network: make(map[uuid.UUID]int), // Assign a unique UUID
		Age: rand.Intn(100), // Randomised age between 0 and 100
		MortalitySalience: false,
		SacrificeChoice: configParam.InitSacrificeChoice,
		ContextSacrifice: "",
	}
}



func (ea *ExtendedAgent) GetName() int {
    return ea.NameID
}

func (ea *ExtendedAgent) SetName(name int) {
    ea.NameID = name
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

func (ea *ExtendedAgent) GetNetwork() map[uuid.UUID]int {
	return ea.Network
}

func (ea *ExtendedAgent) SetNetwork(network map[uuid.UUID]int) {
	ea.Network = network
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

func (ea *ExtendedAgent) GetSacrificeChoice() bool {
    return ea.SacrificeChoice
}


func (ea *ExtendedAgent) GetContextSacrifice() string {
    return ea.ContextSacrifice
}

func (ea *ExtendedAgent) SetContextSacrifice(context string) {
    ea.ContextSacrifice = context
}


// Decision-making logic
func (ea *ExtendedAgent) DecideSacrifice(context string) bool {
    //example will change

	if context == "cause" {
        ea.SacrificeChoice = true
    } else {
        ea.SacrificeChoice = false
    }
    //fmt.Printf("Agent %d decision: %v\n", a.NameID, a.SacrificeChoice)
    ea.ContextSacrifice = context
	fmt.Printf("Agent %d decided to %s for context '%s'\n",
        ea.NameID,
        map[bool]string{true: "sacrifice", false: "not sacrifice"}[ea.SacrificeChoice],
        context)
    return ea.SacrificeChoice
}

func (ea *ExtendedAgent) GetExposedInfo() infra.ExposedAgentInfo {
	return infra.ExposedAgentInfo{
		AgentUUID: ea.GetID(),
	}
}
