package counterServer

import (
	//"TMT/agents"
	"TMT_Attachment/infra"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
)

type TMTServer struct {
	// embed functionality from package...
	// ...and tell BaseServer we're using ICounterAgents
	*server.BaseServer[infra.ITMTAgent]
}

// RunTurn implementation - Count, and then Message

