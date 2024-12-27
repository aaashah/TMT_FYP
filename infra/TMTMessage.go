package infra

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
)

type TMTMessage struct {
	// embed functionality from package
	message.BaseMessage
	// add additional fields
	amountInMessage int
}
