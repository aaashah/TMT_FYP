package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/aaashah/TMT_Attachment/gameRecorder"
	"github.com/google/uuid"
)

type IExtendedAgent interface {
	agent.IAgent[IExtendedAgent]

	//Getters
	GetName() uuid.UUID
	GetAttachment() []float32
	GetNetwork() map[uuid.UUID]float32
	GetAge() int	
	IsMortalitySalient() bool
	GetSelfSacrificeWillingness() float32
	GetPosition() [2]int
	GetWorldviewBinary() uint32
	GetMortality() bool
	//GetContextSacrifice() string
	Move (grid *Grid)

	//Setters
	SetName(name uuid.UUID)
	SetAttachment(attachment []float32)
	SetNetwork(network map[uuid.UUID]float32)
	SetAge(age int)
	SetMortalitySalience(ms bool)
	//SetContextSacrifice(context string)
	
	//DecideSacrifice() bool // Logic for making a self-sacrifice decision.

	//Message functions

	//Info
	GetExposedInfo() ExposedAgentInfo

	// Data Recording
	RecordAgentStatus(instance IExtendedAgent) gameRecorder.AgentRecord

}