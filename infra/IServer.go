package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/google/uuid"
)

type IServer interface {
	agent.IExposedServerFunctions[IExtendedAgent]

	GetAgentByID(agentID uuid.UUID) (IExtendedAgent, bool)
	GetAgentMap() map[uuid.UUID]IExtendedAgent
	//GetAgentPosition(agentID uuid.UUID) ([2]int, bool)

	UpdateAgentRelationship(agentAID, agentBID uuid.UUID, change float32)

	//UpdateAndGetAgentInfo()
	//IsAgentDead(agentID uuid.UUID) bool

}
