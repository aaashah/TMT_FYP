package infra

import "github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"

type IExtendedAgent interface {
	agent.IAgent[IExtendedAgent]

	//Getters
	GetName() int
	GetSacrificeChoice() bool

	//Setters
	SetName(name int)
	DecideSacrifice() bool

	//Message functions

	//Info
	//GetExposedInfo() ExposedAgentInfo

}