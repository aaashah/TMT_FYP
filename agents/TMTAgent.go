package agents

import (
	"TMT_Attachment/infra"
	"fmt"
	"math/rand"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
)

// third tier of composition - embed BaseTMTAgent..
// ... and add 'user specific' fields
type TMTAgent struct {
	*infra.BaseTMTAgent
	amount int
}



// user implementation of DoMessaging ('strategic')
func (uca *TMTAgent) DoMessaging() {
	uca.BroadcastMessage(uca.GetTMTMessage(uca.amount))
	uca.SignalMessagingComplete()
}

// user implementation of Handler ('strategic' - print message to console)
func (uca *TMTAgent) HandleTMTMessage(msg *infra.TMTMessage) {
	fmt.Printf("Sender: %s, Amount: %d\n", msg.GetSender(), msg.GetAmountInMessage())
}

// constructor for UserTMTAgent
func GetUserTMTAgent(funcs agent.IExposedServerFunctions[infra.ITMTAgent]) *TMTAgent {
	return &TMTAgent{
		BaseTMTAgent: infra.GetBaseTMTAgent(funcs),
		amount:           rand.Intn(10),
	}
}
