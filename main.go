package main

//TMTServer "github.com/aaashah/TMT_Attachment/server"

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	//"math/rand"

	baseServer "github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	agents "github.com/aaashah/TMT_Attachment/agents"
	gameRecorder "github.com/aaashah/TMT_Attachment/gameRecorder"
	infra "github.com/aaashah/TMT_Attachment/infra"
	tmtServer "github.com/aaashah/TMT_Attachment/server"
)

// "go run ."
func main() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Create log file with timestamp in name
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFile, err := os.OpenFile("logs/log_"+timestamp+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Create a MultiWriter to write to both the log file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set log output to multiWriter
	log.SetOutput(multiWriter)

	// Remove date and time prefix from log entries
	log.SetFlags(0)

	log.Println("main function started.")
	//rand.Seed(time.Now().UnixNano()) // Seed random number generator
    // Other initialization

	//BEGIN
	// agent configuration:
	agentConfig := agents.AgentConfig{
		InitSacrificeWillingness: 0.2,
	}

	serv := &tmtServer.TMTServer{
		BaseServer: baseServer.CreateBaseServer[infra.IExtendedAgent](
			3, //iterations
			2, //turns per iteration
			50*time.Millisecond, //max duration
			0, //message bandwidth
		),
		ActiveAgents: make(map[uuid.UUID]*agents.ExtendedAgent), // Initialize the activeAgents map
	}
	
	// Set game runner
	serv.SetGameRunner(serv)

	const numAgents int = 5

	// create and initialise agents
	agentPopulation := []infra.IExtendedAgent{}
	grid := infra.CreateGrid(10, 10) // Create a 10x10 grid

	for i := 0; i < numAgents; i++ {
        agentPopulation = append(agentPopulation, agents.CreateExtendedAgents(serv, agentConfig, grid))
    }

	for _, agent := range agentPopulation {
		//agent.SetName(i)
		serv.AddAgent(agent)
		extendedAgent, ok := agent.(*agents.ExtendedAgent)
		if !ok {
			fmt.Printf("Error: agent is not of type *agents.ExtendedAgent\n")
			continue // Skip if type assertion fails
		}

		serv.ActiveAgents[extendedAgent.GetID()] = extendedAgent
		fmt.Printf("Agent %d added with with Age: %d, Attachment: [%.2f, %.2f]\n", agent.GetName(), agent.GetAge(), agent.GetAttachment()[0], agent.GetAttachment()[1])
	}
    
	// Initialize data recorder
	serv.DataRecorder = gameRecorder.CreateRecorder()

	// Start server
	fmt.Println("Starting server")
	serv.Start()

	// custom function to see agent result
	//serv.LogAgentStatus()

	// record data
	serv.DataRecorder.GamePlaybackSummary()
	gameRecorder.ExportToCSV(serv.DataRecorder, "visualisation_output/csv_data")
}