package jetstream

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

var cfg = Config{
	URL: "localhost:4222",
}

func TestKeyValue(t *testing.T) {
	// verifier := token.NewVerifierMock()

	myVal := "some value 1"

	rootCtx := context.Background()

	conn, err := Connect(cfg)
	require.NoError(t, err)
	defer func() {
		conn.conn.Close()
		fmt.Println("conn.conn.Close")
	}()

	kvStore, err := conn.js.CreateKeyValue(rootCtx, jetstream.KeyValueConfig{
		Bucket:      "test-kv",
		Description: "",
	})
	require.NoError(t, err)

	rev, err := kvStore.Put(rootCtx, "somekey1", []byte(myVal))
	require.NoError(t, err)
	fmt.Printf("rev: %v \n", rev)

	asd, err := kvStore.GetRevision(rootCtx, "somekey1", rev)
	require.NoError(t, err)
	fmt.Printf("%#v rev: %v %v \n", asd, asd.Key(), string(asd.Value()))

	keys, err := kvStore.Keys(rootCtx)
	require.NoError(t, err)

	fmt.Printf("keys: %#v", keys)
}

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

	conn, err := Connect(cfg, WithVerifier(verifier))
	require.NoError(t, err)
	defer func() {
		conn.conn.Close()
		fmt.Println("conn.conn.Close")
	}()

	stream, err := conn.CreateStream(rootCtx, "test", "test new stream", []string{"test.*"})
	// stream, err := conn.CreateStream(rootCtx, "test", "test new stream", []string{"test.*"})
	require.NoError(t, err)
	defer func() {
		err = conn.PruneStream(rootCtx, stream.CachedInfo().Config.Name)
		assert.NoError(t, err)
	}()
	fmt.Printf(" %+v \n", stream)

	ctx, span := logger.Span(rootCtx, "main.Method", nil)
	logger.InfoCtx(ctx).Msg("main.Method")

	testPublisher, err := NewPublisher(conn, "test", token.NewTokenGenerator("test-publisher", token.NewSignerMock()))
	require.NoError(t, err)

	consumer1, err := NewConsumer(rootCtx, conn, "test", "group2")
	require.NoError(t, err)

	consunerHandler1 := handlerSvc{t: t, wg: &wg, name: "consumer -- 1"}
	err = consumer1.Consume(ctx, consunerHandler1.handle)
	require.NoError(t, err)

	consumer2, err := NewConsumer(rootCtx, conn, "test", "group2")
	require.NoError(t, err)

	consunerHandler2 := handlerSvc{t: t, wg: &wg, name: "consumer -- 2"}
	err = consumer2.Consume(ctx, consunerHandler2.handle)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	wg.Add(2)
	err = testPublisher.Publish(
		ctx,
		app.NewEventMD(app.WithContentType("type"), app.WithOrderKey("2"), app.WithUniqId("")),
		[]byte("test payload 2222"),
	)
	require.NoError(t, err)

	err = testPublisher.Publish(
		ctx,
		app.NewEventMD(app.WithContentType("type"), app.WithOrderKey("3"), app.WithUniqId("")),
		[]byte("test payload 33333"),
	)
	require.NoError(t, err)

	span.End(nil)
	time.Sleep(1 * time.Second)

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
