package jetstream

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	msgbusv1 "github.com/ggrrrr/btmt-ui/be/common/msgbus/v1"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type ProtoPublisher[T proto.Message] struct {
	publisher *NatsPublisher
	msgMD     protoMetadata
	// jetStream jetstream.Stream
}

// This is only to make sure we have implemented the msgbus.MessagePublisher interface
// anypb.Any is an type  makes compiler happy and it used in the file
var _ (msgbus.MessagePublisher[*anypb.Any]) = (*ProtoPublisher[*anypb.Any])(nil)

func NewCommandPublisher[T proto.Message](
	ctx context.Context,
	conn *NatsConnection,
	msg proto.Message, tokenGenerator token.ServiceTokenGenerator,
) (*ProtoPublisher[T], error) {

	msgMD := createProtoMD(msgTypeCmd, msg)

	publisher, err := NewPublisher(conn, msgMD.publishSubject(), tokenGenerator)
	if err != nil {
		return nil, err
	}

	return &ProtoPublisher[T]{
		publisher: publisher,
		msgMD:     msgMD,
		// jetStream: js,
	}, nil
}

func (p *ProtoPublisher[T]) Publish(ctx context.Context, md msgbus.Metadata, msg T) error {
	if md.Id == uuid.Nil {
		return fmt.Errorf("id is nil")
	}

	payload, err := anypb.New(msg)
	if err != nil {
		return fmt.Errorf("anypb.New %w", err)
	}
	cmd := &msgbusv1.Message{
		Id:          md.Id[:],
		MessageType: p.msgMD.messageType,
		Domain:      p.msgMD.domain,
		Name:        p.msgMD.name,
		Payload:     payload,
		CreatedAt:   timestamppb.Now(),
	}

	bytes, err := proto.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("proto.marshal %w", err)
	}

	return p.publisher.Publish(ctx, md, bytes)

}

func (p *ProtoPublisher[T]) Shutdown() {
	logger.Info().
		Str("Subject", p.msgMD.publishSubject()).
		Msg("Shutdown")
	p.publisher.Shutdown()
}
