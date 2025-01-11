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

	//private
	Attachment []float32 // Attachment orientations: [anxiety, avoidance].
	Kins uuid.UUID
	Heroism float64
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
		Attachment: []float32{0.0, 0.0},
		Kins: uuid.Nil,
		Heroism: 0.0,
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


func (ea *ExtendedAgent) GetContextSacrifice() string {
    return ea.ContextSacrifice
}

func (ea *ExtendedAgent) SetContextSacrifice(context string) {
    ea.ContextSacrifice = context
}


// Decision-making logic
func (ea *ExtendedAgent) DecideSacrifice(context string) bool {
    //example will change

	//anxiety and avoidance influence
	anxietyInfluence := 0.4
	avoidanceInfluence := -0.3

	//combined influence
	attachmentScore := float64(ea.Attachment[0])*anxietyInfluence + float64(ea.Attachment[1])*avoidanceInfluence

	//decision logic
	if ea.MortalitySalience && (ea.Heroism + attachmentScore) > 0.5 {
		ea.SacrificeChoice = true
	} else {
		ea.SacrificeChoice = false
	}
	return ea.SacrificeChoice
}

func (ea *ExtendedAgent) GetExposedInfo() infra.ExposedAgentInfo {
	return infra.ExposedAgentInfo{
		AgentUUID: ea.GetID(),
	}
}
