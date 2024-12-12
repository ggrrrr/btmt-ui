package jetstream

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	msgbusv1 "github.com/ggrrrr/btmt-ui/be/common/msgbus/v1"
)

type msgType string

const msgTypeCmd msgType = "command"
const msgTypeEvent msgType = "event"

type protoMetadata struct {
	// event or command string which is used to create the subject
	// exmaple: event.somedomain.v1.PersonUpdated
	msgType msgType
	// Event or command used for creating the actual proto message
	messageType msgbusv1.MessageType
	// This is string represantaiton of domain and version of the message
	// somedomain.v1
	domain string
	// Name of the command/event example: CreateTrack
	name string
}

func createProtoMD(msgType msgType, msg proto.Message) protoMetadata {
	messageType := msgbusv1.MessageType_MESSAGE_TYPE_UNSPECIFIED
	if msgType == msgTypeEvent {
		messageType = msgbusv1.MessageType_MESSAGE_TYPE_EVENT
	}
	if msgType == msgTypeCmd {
		messageType = msgbusv1.MessageType_MESSAGE_TYPE_COMMAND
	}

	domain := string(msg.ProtoReflect().Descriptor().ParentFile().FullName())
	name := string(msg.ProtoReflect().Descriptor().Name())

	return protoMetadata{
		msgType:     msgType,
		messageType: messageType,
		domain:      domain,
		name:        name,
	}
}

func (m protoMetadata) streamName() string {
	return fmt.Sprintf("app-%s-%s-%s", m.msgType, strings.ReplaceAll(m.domain, ".", "-"), m.name)
}

func (m protoMetadata) consumerSubjects() []string {
	return []string{
		fmt.Sprintf("%s.%s.%s", m.msgType, m.domain, m.name),
		fmt.Sprintf("%s.%s.%s.>", m.msgType, m.domain, m.name),
	}
}

func (m protoMetadata) publishSubject() string {
	return fmt.Sprintf("%s.%s.%s", m.msgType, m.domain, m.name)
}
