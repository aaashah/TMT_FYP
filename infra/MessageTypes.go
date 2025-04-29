package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
)

type WellbeingCheckMessage struct {
	message.BaseMessage
}

type ReplyMessage struct {
	message.BaseMessage
}

func (msg *WellbeingCheckMessage) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleWellbeingCheckMessage(msg)
}

func (msg *ReplyMessage) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleReplyMessage(msg)
}
