package gameRecorder

// "SOMAS_Extended/common"

// AgentRecord is a record of an agent's state at a given turn
type InfraRecord struct {
	// basic info fields
	TurnNumber      int
	IterationNumber int

	Threshold              int  // current threshold set by server
	ThresholdAppliedInTurn bool // whether the threshold was applied in the current turn
}

func NewInfraRecord(turnNumber int, iterationNumber int) InfraRecord {
	return InfraRecord{
		TurnNumber:             turnNumber,
		IterationNumber:        iterationNumber,
	}
}
