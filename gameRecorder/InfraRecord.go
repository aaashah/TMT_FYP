package gameRecorder

// "TMT_Attachment/infra"

// AgentRecord is a record of an agent's state at a given turn
type InfraRecord struct {
	// basic info fields
	TurnNumber      int
	IterationNumber int

	AgentPositions map[[2]int]bool // Stores occupied agent positions (x, y)
	Tombstones     map[[2]int]bool // Stores tombstone locations (x, y)
}

func NewInfraRecord(turnNumber int, iterationNumber int) InfraRecord {
	return InfraRecord{
		TurnNumber:      turnNumber,
		IterationNumber: iterationNumber,
		AgentPositions:  make(map[[2]int]bool),
		Tombstones:      make(map[[2]int]bool),
	}
}
