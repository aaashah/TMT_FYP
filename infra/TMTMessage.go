package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
)

type TMTMessage struct {
	// embed functionality from package
	message.BaseMessage
	// add additional fields
	amountInMessage int
}

// getter for private field
func (cm *TMTMessage) GetAmountInMessage() int {
	return cm.amountInMessage
}

// override of InvokeHandler (visitor pattern)
func (cm *TMTMessage) InvokeMessageHandler(agent ITMTAgent) {
	agent.HandleCounterMessage(cm)
}
