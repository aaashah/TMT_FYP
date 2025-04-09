package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	//"github.com/google/uuid"
)

type IExtendedMessage interface {
	message.IMessage[IExtendedAgent]
}
