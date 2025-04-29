package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
)

type IExtendedMessage interface {
	message.IMessage[IExtendedAgent]
}
