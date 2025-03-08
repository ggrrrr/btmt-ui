package bus

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

type mockApp struct {
	t   *testing.T
	wg  *sync.WaitGroup
	err error
	msg *emailpbv1.EmailMessage
}

// SendEmail implements senderApp.
func (m *mockApp) SendEmail(ctx context.Context, msg *emailpbv1.EmailMessage) error {
	fmt.Printf("\t\tmock.SendEmail88 %+v\n", msg)
	require.NotNil(m.t, msg)
	assert.Equal(m.t, m.msg.FromAccount.Email, msg.FromAccount.Email)
	m.wg.Done()
	return m.err
}

var _ (emailSender) = (*mockApp)(nil)

func T1estPublish(t *testing.T) {
	var err error
	ctx := context.Background()

	err = tracer.ConfigureForTest()
	require.NoError(t, err)
	defer func() {
		tracer.Shutdown(context.Background())
		fmt.Println("logger.Shutdown ;)")
	}()

	conn, err := jetstream.ConnectForTest()
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
	ctx, span := tracer.Tracer("asd").Span(ctx, "testPubklisher")
	err = commandPublisher.Publish(ctx, msgbus.Metadata{Id: testId1}, testEmail1)
	require.NoError(t, err)
	span.End(err)

}

func TestServer(t *testing.T) {
	var err error
	ctx := context.Background()

	err = tracer.ConfigureForTest()
	require.NoError(t, err)
	defer func() {
		tracer.Shutdown(context.Background())
		fmt.Println("loagger.Shutdown ;)")
	}()

	testApp := &mockApp{
		t:  t,
		wg: &sync.WaitGroup{},
	}

	conn, err := jetstream.ConnectForTest()
	require.NoError(t, err)

	consumer, err := jetstream.NewCommandConsumer(ctx, "svc-email", conn, &emailpbv1.SendEmail{},
		func(_ string) *emailpbv1.SendEmail { return new(emailpbv1.SendEmail) },
	)
	require.NoError(t, err)
	defer func() {
		fmt.Println("conn.conn.Close")
		consumer.Shutdown()
	}()

	consumer.Purge(ctx)

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
	testApp.msg = testEmail1.Message

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
