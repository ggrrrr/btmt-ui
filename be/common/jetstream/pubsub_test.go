package jetstream

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

var natsUrl = "localhost:4222"

func TestPublish(t *testing.T) {
	verifier := token.NewVerifierMock()
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

	conn, err := connect("localhost:4222")
	require.NoError(t, err)
	defer func() {
		conn.conn.Close()
		fmt.Println("conn.conn.Close")
	}()

	stream, err := conn.CreateStream(rootCtx, "test", "test stream", []string{"test.*"})
	require.NoError(t, err)
	// defer conn.PruneStream(rootCtx, stream.CachedInfo().Config.Name)
	fmt.Printf(" %+v \n", stream)

	ctx, span := logger.Span(rootCtx, "main.Method", nil)
	logger.InfoCtx(ctx).Msg("main.Method")

	testPublisher, err := NewPublisher(natsUrl, "test", token.NewTokenGenerator("test-publisher", token.NewSignerMock()))
	require.NoError(t, err)
	defer func() {
		_ = testPublisher.Shutdown()
		fmt.Println("testPublisher.Shutdown")
	}()

	consumer, err := NewConsumer(rootCtx, natsUrl, "test", "group2", verifier)
	require.NoError(t, err)
	defer func() {
		_ = consumer.Shutdown()
		fmt.Println("consumer.Shutdown")
	}()

	err = testPublisher.Publish(ctx, "2", []byte("test payload 1asd"))
	require.NoError(t, err)
	wg.Add(1)

	err = consumer.ConsumerLoop(
		func(ctx context.Context, subject string, data []byte) {
			defer wg.Done()
			_, span := logger.Span(ctx, "ConsumerLoop.handler", nil)
			logger.InfoCtx(ctx).Any("data", string(data)).Msg("ConsumerLoop")
			span.End(nil)
		},
	)
	require.NoError(t, err)

	span.End(nil)

	wg.Wait()

}
