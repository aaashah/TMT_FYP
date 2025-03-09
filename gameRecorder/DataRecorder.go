package gameRecorder

import (
	"log"
	"sort"
)

// --------- General External Functions ---------
func Log(message string) {
	log.Println(message)
}

// ✅ Now TurnRecord explicitly encapsulates an Iteration
type TurnRecord struct {
	IterationNumber int  // ✅ Ensures iteration is always before turn
	TurnNumber      int  // ✅ Turns belong inside an iteration
	AgentRecords    []AgentRecord
	InfraRecord     InfraRecord
}

// ✅ Updated Constructor: Iteration comes first
func NewTurnRecord(iterationNumber int, turnNumber int) TurnRecord {
	return TurnRecord{
		IterationNumber: iterationNumber,
		TurnNumber:      turnNumber,
	}
}

// --------- Server Recording Functions ---------
type ServerDataRecorder struct {
	TurnRecords      []TurnRecord // ✅ Stores all iterations & turns
	currentIteration int
	currentTurn      int
}

// ✅ Ensures we can retrieve the current Turn Record
func (sdr *ServerDataRecorder) GetCurrentTurnRecord() *TurnRecord {
	return &sdr.TurnRecords[len(sdr.TurnRecords)-1]
}

// ✅ Creates new recorder
func CreateRecorder() *ServerDataRecorder {
	return &ServerDataRecorder{
		TurnRecords:      []TurnRecord{},
		currentIteration: -1, // Start from -1 so first call goes to 0
		currentTurn:      -1,
	}
}

// ✅ **Properly initializes a new iteration** and resets turns to 0
func (sdr *ServerDataRecorder) RecordNewIteration() {
	sdr.currentIteration += 1
	sdr.currentTurn = -1 // ✅ So that Turn 0 starts fresh in `RecordNewTurn()`
}

// ✅ **Records a new turn within the current iteration**
func (sdr *ServerDataRecorder) RecordNewTurn(agentRecords []AgentRecord, infraRecord InfraRecord) {
    // ✅ Increment turn *before* creating the record, ensuring correct order
    sdr.currentTurn += 1 

    // ✅ Only create Turn 0 the first time it happens in an iteration
    if sdr.currentTurn == 0 {
        sdr.TurnRecords = append(sdr.TurnRecords, NewTurnRecord(sdr.currentIteration, sdr.currentTurn))
        sdr.TurnRecords[len(sdr.TurnRecords)-1].InfraRecord = infraRecord
    }

    sdr.TurnRecords = append(sdr.TurnRecords, NewTurnRecord(sdr.currentIteration, sdr.currentTurn))
    sdr.TurnRecords[len(sdr.TurnRecords)-1].AgentRecords = agentRecords
    sdr.TurnRecords[len(sdr.TurnRecords)-1].InfraRecord = infraRecord
}

// ✅ **Fixes Iteration & Turn Order in CSV & Logs**
func (sdr *ServerDataRecorder) GamePlaybackSummary() {
	log.Printf("\n\nGamePlaybackSummary - playing %v turn records\n", len(sdr.TurnRecords))

	for _, turnRecord := range sdr.TurnRecords {
		log.Printf("\nIteration %v, Turn %v:\n", turnRecord.IterationNumber, turnRecord.TurnNumber)

		// ✅ Print the grid visualization in the logs
		// PrintGrid(turnRecord)

		// ✅ Sort agent records **by ID** for consistent ordering
		sort.Slice(turnRecord.AgentRecords, func(i, j int) bool {
			return turnRecord.AgentRecords[i].AgentID.String() < turnRecord.AgentRecords[j].AgentID.String()
		})

		// ✅ Print agent info correctly
		for _, agentRecord := range turnRecord.AgentRecords {
			log.Printf("Agent %v: ", agentRecord.AgentID)
			agentRecord.DebugPrint()
		}
	}

	// ✅ Creates the HTML visualization
	CreatePlaybackHTML(sdr)
}