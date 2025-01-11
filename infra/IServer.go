package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	//"github.com/google/uuid"
	//"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
)

type IServer interface {
	agent.IExposedServerFunctions[IExtendedAgent]

	//UpdateAndGetAgentInfo()
	//IsAgentDead(agentID uuid.UUID) bool

}