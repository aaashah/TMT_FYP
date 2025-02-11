package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/google/uuid"
)

type IExtendedAgent interface {
	agent.IAgent[IExtendedAgent]

	//Getters
	GetName() int
	GetAttachment() []float32
	GetNetwork() map[uuid.UUID]int
	GetAge() int	
	IsMortalitySalient() bool
	GetSacrificeChoice() bool
	GetContextSacrifice() string

	//Setters
	SetName(name int)
	SetAttachment(attachment []float32)
	SetNetwork(network map[uuid.UUID]int)
	SetAge(age int)
	SetMortalitySalience(ms bool)
	SetContextSacrifice(context string)
	
	DecideSacrifice(context string) bool // Logic for making a self-sacrifice decision.

	//Message functions

	//Info
	GetExposedInfo() ExposedAgentInfo

}