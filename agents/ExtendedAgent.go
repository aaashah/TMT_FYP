package agents

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
)


type ExtendedAgent struct {
	*agent.BaseAgent[infra.IExtendedAgent]
	Server infra.IServer
	NameID int
	Attachment []float32 // Attachment orientations: [anxiety, avoidance].
	Kins uuid.UUID
	Heroism float64
	MortalitySalience bool
	SacrificeChoice bool
	ContextSacrifice string
}

type Attachment struct {
	Anxiety float32
	Avoidance float32
}

func CreateExtendedAgent(server infra.IServer) *ExtendedAgent {
	return &ExtendedAgent{
		BaseAgent: agent.CreateBaseAgent[infra.IExtendedAgent](server),
		Server: server,
		NameID: 0,
		Attachment: []float32{0.5, 0.5},
		Kins: uuid.New(),
		Heroism: 0.5,
		MortalitySalience: false,
		SacrificeChoice: false,
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

func (ea *ExtendedAgent) GetKins() uuid.UUID {
    return ea.Kins
}

func (ea *ExtendedAgent) SetKins(kins uuid.UUID) {
    ea.Kins = kins
}

func (ea *ExtendedAgent) GetHeroism() float64 {
    return ea.Heroism
}

func (ea *ExtendedAgent) SetHeroism(heroism float64) {
    ea.Heroism = heroism
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

func (ea *ExtendedAgent) SetSacrificeChoice(choice bool) {
    ea.SacrificeChoice = choice
}

func (ea *ExtendedAgent) GetContextSacrifice() string {
    return ea.ContextSacrifice
}

func (ea *ExtendedAgent) SetContextSacrifice(context string) {
    ea.ContextSacrifice = context
}

// Decision-making logic
func (ea *ExtendedAgent) DecideSacrifice(context string) bool {
    // Example logic based on attachment and context
    if context == "cause" && ea.MortalitySalience && ea.Attachment[0] > 0.7 {
        return true
    }
    if context == "companion" && ea.MortalitySalience && ea.Attachment[1] < 0.3 {
        return true
    }
    return false
}