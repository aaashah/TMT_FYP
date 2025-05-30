package gameRecorder

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aaashah/TMT_FYP/config"
	"github.com/google/uuid"
)

type Position struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type JSONAgentRecord struct {
	ID                  string   `json:"ID"`
	IsAlive             bool     `json:"IsAlive"`
	Age                 int      `json:"Age"`
	AttachmentStyle     string   `json:"AttachmentStyle"`
	AttachmentAnxiety   float32  `json:"AttachmentAnxiety"`
	AttachmentAvoidance float32  `json:"AttachmentAvoidance"`
	ClusterID           int      `json:"ClusterID"`
	Position            Position `json:"Position"`
	Worldview           uint32   `json:"Worldview"`
	Heroism             int      `json:"Heroism"`
	//MortalitySalience      float32           `json:"MortalitySalience"`
	//WorldviewValidation    float32           `json:"WorldviewValidation"`
	//RelationshipValidation float32           `json:"RelationshipValidation"`
	//ASPDecison		       string            `json:"ASPDecision"`
}

type TurnJSONRecord struct {
	Turn                      int               `json:"TurnNumber"`
	Agents                    []JSONAgentRecord `json:"Agents"`
	EliminatedAgents          []string          `json:"EliminatedAgents"`
	SelfSacrificedAgents      []string          `json:"EliminatedBySelfSacrifice"`
	TotalVolunteers           int               `json:"NumVolunteers"`
	TotalRequiredEliminations int               `json:"TotalRequiredEliminations"`
	TombstoneLocations        []Position        `json:"TombstoneLocations"`
	TempleLocations           []Position        `json:"TempleLocations"`
}

type IterationJSONRecord struct {
	Iteration      int                   `json:"Iteration"`
	Turns          []TurnJSONRecord      `json:"Turns"`
	Thresholds     map[uuid.UUID]float64 `json:"AgentThresholds"`
	NumberOfAgents int                   `json:"NumberOfAgents"`
}

type GameJSONRecord struct {
	Config     config.Config         `json:"Config"`
	Iterations []IterationJSONRecord `json:"Iterations"`
}

func (gjr *GameJSONRecord) AddIteration(record IterationJSONRecord) {
	gjr.Iterations = append(gjr.Iterations, record)
}

func MakeGameRecord(config config.Config) *GameJSONRecord {
	return &GameJSONRecord{
		Config:     config,
		Iterations: make([]IterationJSONRecord, 0),
	}
}

func WriteJSONLog(outputDir string, record *GameJSONRecord) error {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	checkForNaN("GameJSON", record)

	fileName := fmt.Sprintf("%s/output.json", outputDir)
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling iteration JSON: %w", err)
	}

	return os.WriteFile(fileName, data, 0644)
}

func UUIDsToStrings(ids []uuid.UUID) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = id.String()
	}
	return result
}
