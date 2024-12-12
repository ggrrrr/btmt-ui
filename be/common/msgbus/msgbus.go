package msgbus

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type (
	// Handler all messages for a cunsumer ( this could be multiple topics)
	// unmarshaling payload must be in it
	MessageHandlerFunc[T proto.Message] func(ctx context.Context, topic string, md Metadata, msg T) error

	MessageConsumer[T proto.Message] interface {
		Consumer(ctx context.Context, handler MessageHandlerFunc[T])
	}
)

type MessagePublisher[T proto.Message] interface {
	Publish(ctx context.Context, md Metadata, msg T) error
}
