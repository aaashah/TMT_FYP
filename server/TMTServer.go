package server

import (
	"math/rand"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"

	infra "github.com/aaashah/TMT_Attachment/infra"
) 

type TMTServer struct {
	*server.BaseServer[infra.IExtendedAgent]

	agentInfoList []infra.IExtendedAgent

	//server internal state
	//turn int
	//iteration int
}

func int() {
	rand.Seed(time.Now().UnixNano())
}