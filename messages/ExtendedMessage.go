package messages

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	infra "github.com/aaashah/TMT_Attachment/infra"
	//"github.com/google/uuid"
)


type ExtendedMessage struct {
	// embed functionality from package
	message.BaseMessage
	// SenderID uuid.UUID
	// ReceiverID uuid.UUID
}

// func (msg *ExtendedMessage) GetSenderID() uuid.UUID {
// 	return msg.GetSender()
// }

func (msg *ExtendedMessage) InvokeMessageHandler(ea infra.IExtendedAgent) {}