package gameRecorder

import (
	"encoding/json"
	"fmt"
	"os"

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
	NumberOfAgents            int               `json:"NumberOfAgents"`
	EliminatedAgents          []string          `json:"EliminatedAgents"`
	SelfSacrificedAgents      []string          `json:"EliminatedBySelfSacrifice"`
	TotalRequiredEliminations int               `json:"TotalRequiredEliminations"`
	TombstoneLocations        []Position        `json:"TombstoneLocations"`
	TempleLocations           []Position        `json:"TempleLocations"`
}

type IterationJSONRecord struct {
	Iteration int              `json:"Iteration"`
	Turns     []TurnJSONRecord `json:"Turns"`
}

type ConfigJSONRecord struct {
	ProportionAgentsNeeded float32 `json:"PopulationRho"`
}

type GameJSONRecord struct {
	Config     ConfigJSONRecord      `json:"Config"`
	Iterations []IterationJSONRecord `json:"Iterations"`
}

func (gjr *GameJSONRecord) AddIteration(record IterationJSONRecord) {
	gjr.Iterations = append(gjr.Iterations, record)
}

func MakeGameRecord(config ConfigJSONRecord) *GameJSONRecord {
	return &GameJSONRecord{
		Config:     config,
		Iterations: make([]IterationJSONRecord, 0),
	}
}

// func WriteConfigJSONRecord(outputDir string, record ConfigJSONRecord) error {
// 	err := os.MkdirAll(outputDir, os.ModePerm)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output directory: %w", err)
// 	}

// 	fileName := fmt.Sprintf("%s/turn_%03d.json", outputDir, record.Turn)
// 	data, err := json.MarshalIndent(record, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("error marshalling JSON: %w", err)
// 	}
// 	return os.WriteFile(fileName, data, 0644)
// }

// func WriteTurnJSONRecord(outputDir string, record TurnJSONRecord) error {
// 	// Create directory if it doesn't exist
// 	err := os.MkdirAll(outputDir, os.ModePerm)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output directory: %w", err)
// 	}

// 	fileName := fmt.Sprintf("%s/turn_%03d.json", outputDir, record.Turn)
// 	data, err := json.MarshalIndent(record, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("error marshalling JSON: %w", err)
// 	}
// 	return os.WriteFile(fileName, data, 0644)
// }

// func WriteIterationJSONLog(outputDir string, record IterationJSONRecord) error {
// 	err := os.MkdirAll(outputDir, os.ModePerm)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output directory: %w", err)
// 	}

// 	fileName := fmt.Sprintf("%s/iteration_%03d.json", outputDir, record.Iteration)
// 	data, err := json.MarshalIndent(record, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("error marshalling iteration JSON: %w", err)
// 	}

// 	return os.WriteFile(fileName, data, 0644)
// }

func WriteJSONLog(outputDir string, record *GameJSONRecord) error {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

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
