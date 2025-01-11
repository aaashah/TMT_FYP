package server

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"

	agents "github.com/aaashah/TMT_Attachment/agents"
	infra "github.com/aaashah/TMT_Attachment/infra"
	"github.com/google/uuid"
) 

type TMTServer struct {
	*server.BaseServer[infra.IExtendedAgent]

	//agentInfoList []infra.IExtendedAgent
	//mu     sync.Mutex
	context string
	activeAgents map[uuid.UUID]*agents.ExtendedAgent
	// data recorder
	//DataRecorder *gameRecorder.ServerDataRecorder

	//server internal state
	turn int
	//iteration int
	//allAgentsDead bool
	//gameRunner infra.GameRunner
}

func init () {
	rand.Seed(time.Now().UnixNano())
}

func (tserv *TMTServer) RunTurn(i, j int) {
	log.Printf("\n\nIteration %v, Turn %v, current agent count: %v\n", i, j, len(tserv.GetAgentMap()))
	tserv.turn = j
	//1. Agents choose
	for _, agentID := range tserv.activeAgents {
		agentID.DecideSacrifice(tserv.context)
	}

	//2. Eliminate Agents
	remainingAgents := []*agents.ExtendedAgent{}
	for _, agentID := range tserv.activeAgents {
		if !agentID.SacrificeChoice {
			remainingAgents = append(remainingAgents, agentID)
		} else {
			log.Printf("Agent %d has chosen to sacrifice\n", agentID.GetName())
		}
	}
	newActiveAgents := make(map[uuid.UUID]*agents.ExtendedAgent)
	for _, agent := range remainingAgents {
		newActiveAgents[agent.GetID()] = agent
	}
	tserv.activeAgents = newActiveAgents
	fmt.Printf("Turn %d: Ending with %d agents\n", tserv.turn, len(tserv.activeAgents))
	tserv.turn++
}
