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
	GetAttachment() Attachment
	GetNetwork() map[uuid.UUID]float32
	GetAge() int
	//IsMortalitySalient() bool
	GetSelfSacrificeWillingness() float64
	GetPosition() PositionVector
	SetPosition(PositionVector)
	GetWorldviewBinary() uint32
	GetMortality() bool
	//GetContextSacrifice() string
	// Move(grid *Grid)
	GetTargetPosition(grid *Grid) (PositionVector, bool)
	GetClusterID() int
	UpdateRelationship(agentID uuid.UUID, change float32)
	GetASPDecision(grid *Grid) int
	IncrementClusterEliminations(n int)
	IncrementNetworkEliminations(n int)
	IncrementHeroism()

	//Setters
	// SetName(name uuid.UUID)
	// SetAttachment(attachment []float32)
	// SetNetwork(network map[uuid.UUID]float32)
	IncrementAge()
	//SetMortalitySalience(ms bool)
	SetClusterID(id int)
	//SetContextSacrifice(context string)

	//Message functions

	//Info
	GetExposedInfo() ExposedAgentInfo
	AgentInitialised()

	// Data Recording
	RecordAgentStatus(instance IExtendedAgent) gameRecorder.AgentRecord

	// Updaters
	UpdateSocialNetwork(uuid.UUID, float32)
}
