package msgbusv1

import (
	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

var _ (logger.TraceDataExtractor) = (*Message)(nil)

func (x *Message) Extract() logger.TraceData {
	id, _ := uuid.ParseBytes(x.Id)
	return logger.TraceData{
		"message.id":               logger.TraceValueString(id.String()),
		"message.domain":           logger.TraceValueString(x.Domain),
		"message.name":             logger.TraceValueString(x.Name),
		"message.type":             logger.TraceValueString(x.MessageType.String()),
		"message.payload.type.url": logger.TraceValueString(x.Payload.TypeUrl),
	}
}
