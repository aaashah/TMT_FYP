package infra

import (
	//"github.com/google/uuid"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
)

type WellbeingCheckMessage struct {
	message.BaseMessage
	//WellbeingMessage string
}

type ReplyMessage struct {
	message.BaseMessage
	//AckMessage string
}

func (msg *WellbeingCheckMessage) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleWellbeingCheckMessage(msg)
}

func (msg *ReplyMessage) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleReplyMessage(msg)
}
