package gameRecorder

import (
	"log"

	"github.com/google/uuid"
)

type AgentRecord struct {
	// basic info fields
	TurnNumber      int
	IterationNumber int
	AgentID         uuid.UUID
	AgentAge 		    int


	// turn-specific fields
	IsAlive            bool
	PositionX                int
	PositionY                int
	SelfSacrificeWillingness float32
	AttachmentAnxiety        float32
	AttachmentAvoidance      float32
	Worldview                uint32
	//WorldviewValidation      float64
	//RelationshipValidation   float64
	
	// special indicator fields for agents
	Died bool // of natural causes
	SpecialNote string
}

func NewAgentRecord(agentID uuid.UUID, agentAge int, positionX int, positionY int, sacrificeWillingness float32, attachmentAnxiety float32, attachmentAvoidance float32, worldview uint32, specialNote string) AgentRecord {
	return AgentRecord{
		AgentID:            agentID,
		AgentAge: 			agentAge,
		PositionX:          positionX,
		PositionY:          positionY,
		SelfSacrificeWillingness: sacrificeWillingness,
		AttachmentAnxiety:        attachmentAnxiety,
		AttachmentAvoidance:      attachmentAvoidance,
		Worldview:	worldview,
		//WorldviewValidation:      worldviewValidation,
		//RelationshipValidation:   relationshipValidation,
		//Died:			   false,
		SpecialNote:        specialNote,
	}
}

func (ar *AgentRecord) DebugPrint() {
	// log.Printf("Agent ID: %v\n", ar.AgentID)
	if !ar.IsAlive {
		log.Printf("[DEAD] ")
	}
	log.Printf("Agent ID: %v\n", ar.AgentID)
	// log.Printf("Agent Contribution: %v\n", ar.agent.GetActualContribution(ar.agent))
	// log.Printf("Agent Stated Contribution: %v\n", ar.agent.GetStatedContribution(ar.agent))
	// log.Printf("Agent Withdrawal: %v\n", ar.agent.GetActualWithdrawal(ar.agent))
	// log.Printf("Agent Stated Withdrawal: %v\n", ar.agent.GetStatedWithdrawal(ar.agent))
}