package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/google/uuid"
)

type GameRunner interface {
	RunStartOfIteration(iteration int)
	RunTurn(iteration int, turn int)
	RunEndOfIteration(iteration int)
}

type IAgentOperations[T any] interface {
	GetAgentMap() map[uuid.UUID]T
	AddAgent(agentToAdd T)
	RemoveAgent(agentToRemove T)
}

type IServer interface {
	agent.IExposedServerFunctions[IExtendedAgent]

	GetIterations() int
	GetTurns() int
	SetGameRunner(gameRunner GameRunner)
	Start()
}