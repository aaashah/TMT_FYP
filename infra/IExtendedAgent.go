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
	GetSelfSacrificeWillingness() float32
	GetPosition() PositionVector
	SetPosition(PositionVector)
	GetWorldviewBinary() uint32
	//GetMortality() bool
	GetTelomere() float32
	IsAlive() bool
	//GetContextSacrifice() string
	// Move(grid *Grid)
	GetTargetPosition(grid *Grid) (PositionVector, bool)
	GetClusterID() int
	UpdateRelationship(agentID uuid.UUID, change float32)
	RemoveRelationship(agentID uuid.UUID)
	GetASPDecision(grid *Grid) ASPDecison
	GetPTSParams() PTSParams
	IncrementClusterEliminations(n int)
	IncrementNetworkEliminations(n int)
	IncrementHeroism()
	GetHeroism() int

	//Setters
	// SetName(name uuid.UUID)
	// SetAttachment(attachment []float32)
	// SetNetwork(network map[uuid.UUID]float32)
	IncrementAge()
	//SetMortalitySalience(ms bool)
	SetClusterID(id int)
	//SetContextSacrifice(context string)
	MarkAsDead()
	UpdateEsteem(id uuid.UUID, isCheck bool)
	SetWorldviewBinary(worldview uint32)
	SetParents(parent1, parent2 uuid.UUID)
	AddDescendant(descendant uuid.UUID)

	//Message functions
	CreateWellbeingCheckMessage() *WellbeingCheckMessage
	CreateReplyMessage() *ReplyMessage
	HandleWellbeingCheckMessage(msg *WellbeingCheckMessage)
	HandleReplyMessage(msg *ReplyMessage)

	//Info
	GetExposedInfo() ExposedAgentInfo
	AgentInitialised()

	// Data Recording
	RecordAgentStatus(instance IExtendedAgent) gameRecorder.AgentRecord

	// Updaters
	UpdateSocialNetwork(uuid.UUID, float32)
}
