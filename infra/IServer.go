package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	//"github.com/google/uuid"
)

type IServer interface {
	agent.IExposedServerFunctions[IExtendedAgent]

}