package eda

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

type mockApp struct {
	err error
}

// SendEmail implements senderApp.
func (m *mockApp) SendEmail(ctx context.Context, msg *emailpbv1.EmailMessage) error {
	fmt.Printf("msg %#v\n", msg)
	return m.err
}

var cfg = jetstream.Config{
	URL: "localhost:4222",
}
var _ (senderApp) = (*mockApp)(nil)

func TestServer(t *testing.T) {
	verifier := token.NewVerifierMock()
	var err error
	ctx := context.Background()

	os.Setenv("OTEL_COLLECTOR", "localhost:4317")
	os.Setenv("SERVICE_NAME", "test-service")

	testApp := &mockApp{}

	err = logger.ConfigureOtel(ctx)
	require.NoError(t, err)
	defer func() {
		logger.Shutdown()
		fmt.Println("logger.Shutdown ;)")
	}()

	conn, err := jetstream.Connect(cfg, jetstream.WithVerifier(verifier))
	require.NoError(t, err)
	defer func() {
		fmt.Println("conn.conn.Close")
	}()

	err = NewCommandHandler(ctx, testApp, conn)
	require.NoError(t, err)

	// command := &eventv1.Command{}

	// sendEmail := &emailpbv1.SendEmail{
	// 	Id: "some_id",
	// 	Message: &emailpbv1.EmailMessage{
	// 		FromAccount: &emailpbv1.SenderAccount{
	// 			Realm: "localhost",
	// 		},
	// 		ToEmail: []*emailpbv1.ToEmail{&emailpbv1.ToEmail{
	// 			Name:  "from",
	// 			Email: "from@me",
	// 		}},
	// 		// Body: "some body",
	// 	},
	// }
	// details, err := ptypes.MarshalAny(sendEmail)
	// command.Details = details
	// // details, err := types.MarshalAny(sendEmail)
	// // err = command.Details.MarshalFrom(sendEmail)
	// require.NoError(t, err)

	// commandBytes, err := proto.Marshal(command)
	// require.NoError(t, err)

	// publisher, err := jetstream.NewPublisher(conn, "svc-email", token.NewTokenGenerator("test-publisher", token.NewSignerMock()))
	// require.NoError(t, err)

	// err = publisher.Publish(ctx, app.NewMetadata(), commandBytes)
	// require.NoError(t, err)

	time.Sleep(10 * time.Second)
}

func TestMsg(t *testing.T) {

	test := &emailpbv1.SendEmail{}
	fmt.Printf("%#v \n", test.ProtoReflect().Descriptor())
	fmt.Printf("%#v \n", test.ProtoReflect().Descriptor().Name())
	fmt.Printf("%#v \n", test.ProtoReflect().Descriptor().ParentFile())

}
