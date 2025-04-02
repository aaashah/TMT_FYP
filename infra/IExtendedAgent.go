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
	GetClusterID() int
	UpdateRelationship(agentID uuid.UUID, change float32)
	DecideSacrifice() float32

	//Setters
	SetName(name uuid.UUID)
	SetAttachment(attachment []float32)
	SetNetwork(network map[uuid.UUID]float32)
	SetAge(age int)
	SetMortalitySalience(ms bool)
	SetClusterID(id int)
	//SetContextSacrifice(context string)


	//Message functions

	//Info
	GetExposedInfo() ExposedAgentInfo
	AgentInitialised()

	// Data Recording
	RecordAgentStatus(instance IExtendedAgent) gameRecorder.AgentRecord

}