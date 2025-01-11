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
	GetKins() uuid.UUID
	GetHeroism() float64
	IsMortalitySalient() bool
	GetSacrificeChoice() bool
	GetContextSacrifice() string

	//Setters
	SetName(name int)
	SetAttachment(attachment []float32)
	SetKins(kins uuid.UUID)
	SetHeroism(heroism float64)
	SetMortalitySalience(ms bool)
	SetContextSacrifice(context string)
	
	DecideSacrifice(context string) bool // Logic for making a self-sacrifice decision.

	//Message functions

	//Info
	GetExposedInfo() ExposedAgentInfo

}