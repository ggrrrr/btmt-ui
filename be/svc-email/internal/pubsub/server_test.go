package eda

import (
	"context"
	"fmt"
	"os"
	"testing"

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

}
