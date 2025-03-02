package jetstream

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	msgbusv1 "github.com/ggrrrr/btmt-ui/be/common/msgbus/v1"
)

type (
	ProtoConsumer[T proto.Message] struct {
		consumer     *NatsConsumer
		jetStream    jetstream.Stream
		internalMsg  protoMetadata
		newValueFunc func(topic string) T
	}
)

// This is only to make sure we have implemented the msgbus.MessageConsumer interface
// anypb.Any is an type  makes compiler happy and it used in the file
var _ (msgbus.MessageConsumer[*anypb.Any]) = (*ProtoConsumer[*anypb.Any])(nil)

func NewCommandConsumer[T proto.Message](
	ctx context.Context,
	consumerGroup string,
	conn *NatsConnection,
	msg T,
	newValueFunc func(topic string) T,
) (*ProtoConsumer[T], error) {

	msgMD := createProtoMD(msgTypeCmd, msg)
	streamName := msgMD.streamName()

	js, err := conn.CreateStream(ctx, streamName, fmt.Sprintf("stream for %s command", msgMD.publishSubject()), msgMD.consumerSubjects())
	if err != nil {
		return nil, fmt.Errorf("conn.CreateStream %w", err)
	}

	consumer, err := NewConsumer(ctx, conn, streamName, consumerGroup)
	if err != nil {
		return nil, err
	}

	return &ProtoConsumer[T]{
		consumer:     consumer,
		jetStream:    js,
		internalMsg:  msgMD,
		newValueFunc: newValueFunc,
	}, nil
}

func (c *ProtoConsumer[T]) Consumer(ctx context.Context, handler msgbus.MessageHandlerFunc[T]) {
	c.consumer.Consume(ctx, func(ctx context.Context, subject string, md msgbus.Metadata, data []byte) error {
		msg := &msgbusv1.Message{}
		err := proto.Unmarshal(data, msg)
		if err != nil {
			return fmt.Errorf("proto.Unmarshal %w", err)
		}

		md.Id = uuid.UUID(msg.Id)

		payload := c.newValueFunc(subject)

		err = anypb.UnmarshalTo(msg.Payload, payload, proto.UnmarshalOptions{})
		if err != nil {
			return fmt.Errorf("anypb.UnmarshalTo %w", err)
		}

		handler(ctx, subject, md, payload)
		return nil
	})
}

func (c *ProtoConsumer[T]) Shutdown() {
	log.Log().Info("Shutdown")
	c.consumer.Shutdown()
}
