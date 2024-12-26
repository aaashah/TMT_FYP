package infra

import "github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"

type ITMTAgent interface {
	// embed functionality from package
	agent.IAgent[ITMTAgent]
	// perfom counting action
	DoCount()
	// getter for count value (agents are injected as interfaces,...
	// ...not structs, so Count is not visible)
	GetCount() int
	// perform messaging action
	DoMessaging()
	// get CounterMessage (QoL - convenient to call from agent)
	GetCounterMessage(int) *TMTMessage
	// handler for CounterMessage (visitor design pattern)
	HandleCounterMessage(*TMTMessage)
}
