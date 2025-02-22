package bus

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

type mockApp struct {
	wg  *sync.WaitGroup
	err error
}

// SendEmail implements senderApp.
func (m *mockApp) SendEmail(ctx context.Context, msg *emailpbv1.EmailMessage) error {
	fmt.Printf("\t\tmock.SendEmail88 %+v\n", msg)
	m.wg.Done()
	return m.err
}

var cfg = jetstream.Config{
	URL: "localhost:4222",
}
var _ (emailSender) = (*mockApp)(nil)

func TestPublish(t *testing.T) {
	verifier := token.NewVerifierMock()
	var err error
	ctx := context.Background()

	os.Setenv("OTEL_COLLECTOR", "localhost:4317")
	os.Setenv("SERVICE_NAME", "test-service")
	err = logger.ConfigureOtel(ctx, "devapp", logger.DevConfig)
	require.NoError(t, err)
	defer func() {
		logger.Shutdown()
		fmt.Println("logger.Shutdown ;)")
	}()

	conn, err := jetstream.Connect(cfg, jetstream.WithVerifier(verifier))
	require.NoError(t, err)

	commandPublisher, err := jetstream.NewCommandPublisher[*emailpbv1.SendEmail](
		ctx,
		conn,
		&emailpbv1.SendEmail{},
		token.NewTokenGenerator("test-publisher", token.NewSignerMock()),
	)
	require.NoError(t, err)

	testId1 := uuid.New()
	testEmail1 := &emailpbv1.SendEmail{
		Id: testId1.String(),
		Message: &emailpbv1.EmailMessage{
			FromAccount: &emailpbv1.SenderAccount{
				Realm: "localhost",
				Name:  "from",
				Email: "from@me",
			},
			ToAddresses: &emailpbv1.ToAddresses{
				ToEmail: []*emailpbv1.EmailAddr{
					&emailpbv1.EmailAddr{Email: "me@email.com"},
				},
			},
			Body: &emailpbv1.EmailMessage_RawBody{
				RawBody: &emailpbv1.RawBody{
					ContentType: "",
					Subject:     "subject 1",
					Body:        "body1",
				},
			},
		},
	}
	ctx, span := logger.Span(ctx, "testPubklisher", nil)
	err = commandPublisher.Publish(ctx, msgbus.Metadata{Id: testId1}, testEmail1)
	require.NoError(t, err)
	span.End(err)

}

func TestServer(t *testing.T) {
	verifier := token.NewVerifierMock()
	var err error
	ctx := context.Background()

	os.Setenv("OTEL_COLLECTOR", "localhost:4317")
	os.Setenv("SERVICE_NAME", "test-service")
	err = logger.ConfigureOtel(ctx, "devapp", logger.DevConfig)
	require.NoError(t, err)
	defer func() {
		logger.Shutdown()
		fmt.Println("loagger.Shutdown ;)")
	}()

	testApp := &mockApp{
		wg: &sync.WaitGroup{},
	}

	conn, err := jetstream.Connect(cfg, jetstream.WithVerifier(verifier))
	require.NoError(t, err)

	consumer, err := jetstream.NewCommandConsumer(ctx, "svc-email", conn, &emailpbv1.SendEmail{},
		func(_ string) *emailpbv1.SendEmail { return new(emailpbv1.SendEmail) },
	)
	require.NoError(t, err)
	defer func() {
		fmt.Println("conn.conn.Close")
		consumer.Shutdown()
	}()

	err = Start(ctx, testApp, consumer)
	require.NoError(t, err)

	commandPublisher, err := jetstream.NewCommandPublisher[*emailpbv1.SendEmail](
		ctx,
		conn,
		&emailpbv1.SendEmail{},
		token.NewTokenGenerator("test-publisher", token.NewSignerMock()),
	)
	require.NoError(t, err)

	defer func() {
		commandPublisher.Shutdown()
	}()

	testId1 := uuid.New()
	testEmail1 := &emailpbv1.SendEmail{
		Id: testId1.String(),
		Message: &emailpbv1.EmailMessage{
			FromAccount: &emailpbv1.SenderAccount{
				Realm: "localhost",
				Name:  "from",
				Email: "from@me",
			},
		},
	}

	err = commandPublisher.Publish(ctx, msgbus.Metadata{Id: testId1}, testEmail1)
	require.NoError(t, err)
	testApp.wg.Add(1)

	testApp.wg.Wait()
}

func TestMsg(t *testing.T) {
	test := &emailpbv1.SendEmail{}

	domain := string(test.ProtoReflect().Descriptor().ParentFile().Package().Parent().Name())
	domain1 := string(test.ProtoReflect().Descriptor().ParentFile().Package().Parent())
	name := string(test.ProtoReflect().Descriptor().Name())

	fmt.Printf("domain: %#v \n", domain)
	fmt.Printf("domain1: %#v \n", domain1)
	fmt.Printf("name: %#v \n", name)

	// fmt.Printf("%#v \n", test.ProtoReflect().Descriptor())
	// fmt.Printf("%#v \n", test.ProtoReflect().Descriptor().Name())
	// fmt.Printf("%#v \n", test.ProtoReflect().Descriptor().ParentFile())

}
