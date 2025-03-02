package msgbusv1

import (
	"log/slog"

	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

var _ (td.TraceDataExtractor) = (*Message)(nil)

func (x *Message) Extract() *td.TraceData {
	id, _ := uuid.ParseBytes(x.Id)
	return &td.TraceData{
		KV: map[string]slog.Value{
			"message.id":               slog.StringValue(id.String()),
			"message.domain":           slog.StringValue(x.Domain),
			"message.name":             slog.StringValue(x.Name),
			"message.type":             slog.StringValue(x.MessageType.String()),
			"message.payload.type.url": slog.StringValue(x.Payload.TypeUrl),
		},
	}
}
