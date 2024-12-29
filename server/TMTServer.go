package environmentServer

import (
	"math/rand"
	"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"

	infra "github.com/aaashah/TMT_Attachment/infra"
) 

type TMTServer struct {
	*server.BaseServer[infra.IExtendedAgent]

	agentInfoList []infra.IExtendedAgent
	mu     sync.Mutex

	// data recorder
	//DataRecorder *gameRecorder.ServerDataRecorder

	//server internal state
	turn int
	iteration int
}

func int() {
	rand.Seed(time.Now().UnixNano())
}