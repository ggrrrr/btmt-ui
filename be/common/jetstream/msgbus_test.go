package jetstream

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	msgbusv1 "github.com/ggrrrr/btmt-ui/be/common/msgbus/v1"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

func TestCreateProtoMD(t *testing.T) {

	tests := []struct {
		name    string
		msgType msgType
		from    proto.Message
		result  protoMetadata
	}{
		{
			name:    "v1",
			msgType: msgTypeCmd,
			from:    &msgbusv1.TestCommand{},
			result: protoMetadata{
				msgType:     msgTypeCmd,
				messageType: msgbusv1.MessageType_MESSAGE_TYPE_COMMAND,
				domain:      "msgbus.v1",
				name:        "TestCommand",
			},
		},
		{
			name:    "app",
			msgType: msgTypeEvent,
			from:    &msgbusv1.TestCommand{},
			result: protoMetadata{
				msgType:     msgTypeEvent,
				messageType: msgbusv1.MessageType_MESSAGE_TYPE_EVENT,
				domain:      "msgbus.v1",
				name:        "TestCommand",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := createProtoMD(tc.msgType, tc.from)
			require.Equal(t, tc.result, actual)
		})
	}

}

func TestPubSubNats(t *testing.T) {
	wg := sync.WaitGroup{}

	var err error
	rootCtx := context.Background()

	os.Setenv("OTEL_COLLECTOR", "localhost:4317")
	os.Setenv("SERVICE_NAME", "test-service")

	err = logger.ConfigureOtel(rootCtx)
	require.NoError(t, err)
	defer func() {
		logger.Shutdown()
		fmt.Println("logger.Shutdown ;)")
	}()

	verifier := token.NewVerifierMock()
	conn, err := Connect(cfg, WithVerifier(verifier))
	require.NoError(t, err)
	defer func() {
		conn.conn.Close()
		fmt.Println("conn.conn.Close")
	}()

	commandPublisher, err := NewCommandPublisher[*msgbusv1.TestCommand](
		rootCtx,
		conn,
		&msgbusv1.TestCommand{},
		token.NewTokenGenerator("test-publisher", token.NewSignerMock()),
	)
	require.NoError(t, err)

	commandConsumer, err := NewCommandConsumer(
		rootCtx, "test-app-1",
		conn, &msgbusv1.TestCommand{},
		func(_ string) *msgbusv1.TestCommand {
			return &msgbusv1.TestCommand{}
		},
	)
	require.NoError(t, err)
	defer func() {
		err = conn.PruneStream(rootCtx, commandConsumer.internalMsg.streamName())
		assert.NoError(t, err)
	}()

	testId := uuid.New()

	err = commandPublisher.Publish(
		rootCtx,
		msgbus.Metadata{Id: testId, RetryCounter: 2, RetryAfter: 4 * time.Second, RetryTopic: "sometopic"},
		&msgbusv1.TestCommand{Id: testId.String(), Name: "test command", Info: "some info"},
	)
	require.NoError(t, err)
	wg.Add(1)

	commandConsumer.Consumer(rootCtx, func(ctx context.Context, topic string, md msgbus.Metadata, msg *msgbusv1.TestCommand) error {
		fmt.Printf("\t\t md: %+v \n", md)
		fmt.Printf("\t\t topic: %v \n", topic)
		fmt.Printf("\t\t msg: %+v \n", msg)
		// fmt.Printf("\t md: %#v \n")

		wg.Done()
		return nil
	})

	wg.Wait()
}
