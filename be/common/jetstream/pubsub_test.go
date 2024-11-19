package jetstream

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
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

	stream, err := conn.CreateStream(rootCtx, "test", "test new stream", []string{"test.*"})
	require.NoError(t, err)
	defer conn.PruneStream(rootCtx, stream.CachedInfo().Config.Name)
	fmt.Printf(" %+v \n", stream)

	ctx, span := logger.Span(rootCtx, "main.Method", nil)
	logger.InfoCtx(ctx).Msg("main.Method")

	testPublisher, err := NewPublisher(natsUrl, "test", token.NewTokenGenerator("test-publisher", token.NewSignerMock()))
	require.NoError(t, err)
	defer func() {
		_ = testPublisher.Shutdown()
		fmt.Println("testPublisher.Shutdown")
	}()

	consumer1, err := NewConsumer(rootCtx, natsUrl, "test", "group2", verifier)
	require.NoError(t, err)
	defer func() {
		_ = consumer1.Shutdown()
		fmt.Println("consumer.Shutdown")
	}()
	consunerHandler1 := handlerSvc{t: t, wg: &wg, name: "consumer -- 1"}
	err = consumer1.ConsumerLoop(consunerHandler1.handle)
	require.NoError(t, err)

	consumer2, err := NewConsumer(rootCtx, natsUrl, "test", "group2", verifier)
	require.NoError(t, err)
	defer func() {
		_ = consumer2.Shutdown()
		fmt.Println("consumer.Shutdown")
	}()
	consunerHandler2 := handlerSvc{t: t, wg: &wg, name: "consumer -- 2"}
	err = consumer2.ConsumerLoop(consunerHandler2.handle)
	require.NoError(t, err)

	time.Sleep(10 * time.Second)

	err = testPublisher.Publish(ctx, "2", []byte("test payload 2222"))
	require.NoError(t, err)
	wg.Add(1)

	err = testPublisher.Publish(ctx, "3", []byte("test payload 33333"))
	require.NoError(t, err)
	wg.Add(1)

	span.End(nil)
	time.Sleep(10 * time.Second)

	wg.Wait()

}

type handlerSvc struct {
	t    *testing.T
	wg   *sync.WaitGroup
	name string
}

func (h handlerSvc) handle(ctx context.Context, subject string, data []byte) {
	defer h.wg.Done()
	ctx, span := logger.Span(ctx, "ConsumerLoop.handler", nil)
	authInfo := roles.AuthInfoFromCtx(ctx)
	assert.Equalf(h.t, "mockuser", authInfo.Subject, "authInfo is not set %#v", authInfo)
	logger.InfoCtx(ctx).Str("group", h.name).Any("data", string(data)).Msg("ConsumerLoop")
	span.End(nil)
}
