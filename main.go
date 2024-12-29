package main

//TMTServer "github.com/aaashah/TMT_Attachment/server"

import (
	"fmt"
	"log"
	"time"

	baseServer "github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	agents "github.com/aaashah/TMT_Attachment/agents"
	infra "github.com/aaashah/TMT_Attachment/infra"
	tmtServer "github.com/aaashah/TMT_Attachment/server"
)

// "go run ."
func main() {
	//TMTServer.StartServer()
	log.Println("main function started.")

	serv := &tmtServer.TMTServer{
		// note: the zero turn is used for team forming
		BaseServer: baseServer.CreateBaseServer[infra.IExtendedAgent](
			3,                   //  iterations
			120,                 //  turns per iteration
			50*time.Millisecond, //  max duration
			10),                 //  message bandwidth
	}
	
	// Create agents
	const numAgents int = 10
	agentPopulation := []infra.IExtendedAgent{}
	for i := 0; i < numAgents; i++ {
		//agentPopulation = append(agentPopulation, agents.Team4_CreateAgent(serv, agentConfig))
		agentPopulation = append(agentPopulation, agents.CreateExtendedAgent(serv))
		// Add other teams' agents here
	}

	

	// Set game runner
	serv.SetGameRunner(serv)

	// Start server
	fmt.Println("Starting server")
	serv.Start()
}
