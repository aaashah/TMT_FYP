package infra

import (
	"fmt"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
)

type BaseTMTAgent struct {
	// embed functionality of package
	*agent.BaseAgent[ITMTAgent]
	// add additional fields for all agents in simulator
	Count int
}

// base implementation of DoCount
func (bca *BaseTMTAgent) DoCount() {
	fmt.Printf("%s is counting...\n", bca.GetID())
	bca.Count += 1
}

// base implementation of DoMessaging (just end straight away)
func (bca *BaseTMTAgent) DoMessaging() {
	bca.SignalMessagingComplete()
}

// 'correct' implementation of GetCount - override not needed
func (bca *BaseTMTAgent) GetCount() int {
	return bca.Count
}

// base implmentation of HandleCounterMessage (just ignore)
func (bca *BaseTMTAgent) HandleCounterMessage(msg *TMTMessage) {}

// constructor for CounterMessage (QoL - callable from agent)
func (bca *BaseTMTAgent) GetCounterMessage(amt int) *TMTMessage {
	return &TMTMessage{
		BaseMessage:     bca.CreateBaseMessage(),
		amountInMessage: amt,
	}
}

// constructor for BaseCounterAgent
func GetBaseCounterAgent(funcs agent.IExposedServerFunctions[ITMTAgent]) *BaseTMTAgent {
	return &BaseTMTAgent{
		BaseAgent: agent.CreateBaseAgent(funcs),
	}
}
