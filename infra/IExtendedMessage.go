package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	"github.com/google/uuid"
)

type IExtendedMessageMessage interface {
	message.IMessage[IExtendedAgent]
	GetSenderID() uuid.UUID
	GetReceiverID() uuid.UUID
}