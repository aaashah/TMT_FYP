package main

//TMTServer "github.com/aaashah/TMT_Attachment/server"

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

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
	// agentConfig := agents.AgentConfig{
	// 	InitSacrificeWillingness: 0.2,
	// }
	grid := infra.NewGrid(70, 30) // Create grid

	serv := &tmtServer.TMTServer{
		BaseServer: baseServer.CreateBaseServer[infra.IExtendedAgent](
			3,                   //iterations
			5,                   //turns per iteration
			50*time.Millisecond, //max duration
			0,                   //message bandwidth
		),
		//ActiveAgents: make(map[uuid.UUID]*agents.ExtendedAgent), // Initialize the activeAgents map
		Grid: grid,
	}

	// Set game runner
	serv.SetGameRunner(serv)

	const numAgents int = 16

	// create and initialise agents
	agentPopulation := []infra.IExtendedAgent{}
	//grid := infra.NewGrid(70, 30) // Create grid

	//funcs:= &IExposedServerFunctions[infra.IExtendedAgent]

	for i := 0; i < numAgents; i += 4 {
		agentPopulation = append(agentPopulation, agents.CreateSecureAgent(serv, grid))
		agentPopulation = append(agentPopulation, agents.CreateDismissiveAgent(serv, grid))
		agentPopulation = append(agentPopulation, agents.CreatePreoccupiedAgent(serv, grid))
		agentPopulation = append(agentPopulation, agents.CreateFearfulAgent(serv, grid))
	}

	// Set probability p for Erdős–Rényi network

	for _, agent := range agentPopulation {
		//agent.SetName(i)
		serv.AddAgent(agent)

		agent.AgentInitialised() // Call the method to print agent details

		//fmt.Printf("Agent %v added with with Age: %d, Attachment: [%.2f, %.2f]\n", agent.GetID(), agent.GetAge(), agent.GetAttachment()[0], agent.GetAttachment()[1])
	}

	//const connectionProbability = 0.3 // Adjust as needed
	// Initialize social network after agents are created
	//serv.InitialiseRandomNetwork(connectionProbability)

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
	//gameRecorder.CreateGridHTML(serv.DataRecorder, "visualisation_output/csv_data")
}
