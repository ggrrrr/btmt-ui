package jetstream

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
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

	testId := uuid.New()

	topic := "test"
	verifier := token.NewVerifierMock()
	wg := sync.WaitGroup{}

	var err error
	rootCtx := context.Background()

	os.Setenv("OTEL_COLLECTOR", "localhost:4317")
	os.Setenv("SERVICE_NAME", "test-service")

	err = tracer.Configure(rootCtx, "testapp", tracer.Config{})
	require.NoError(t, err)
	defer func() {
		tracer.Shutdown(rootCtx)
		fmt.Println("logger.Shutdown ;)")
	}()

	conn, err := Connect(cfg, WithVerifier(verifier))
	require.NoError(t, err)
	defer func() {
		conn.conn.Close()
		fmt.Println("conn.conn.Close")
	}()

	stream, err := conn.CreateStream(rootCtx, "test", "test new stream", []string{"test", "test.*"})
	// stream, err := conn.CreateStream(rootCtx, "test", "test new stream", []string{"test.*"})
	require.NoError(t, err)
	defer func() {
		err = conn.PruneStream(rootCtx, stream.CachedInfo().Config.Name)
		assert.NoError(t, err)
	}()
	fmt.Printf(" %+v \n", stream)

	ctx, span := tracer.Tracer("asd").Span(rootCtx, "main.Method")
	log.Log().InfoCtx(ctx, "main.Method")

	testPublisher, err := NewPublisher(conn, topic, token.NewTokenGenerator("test-publisher", token.NewSignerMock()))
	require.NoError(t, err)

	consumer1, err := NewConsumer(rootCtx, conn, topic, "group2")
	require.NoError(t, err)

	consunerHandler1 := handlerSvc{t: t, wg: &wg, name: "consumer -- 1"}
	err = consumer1.Consume(ctx, consunerHandler1.handle)
	require.NoError(t, err)

	consumer2, err := NewConsumer(rootCtx, conn, topic, "group2")
	require.NoError(t, err)

	consunerHandler2 := handlerSvc{t: t, wg: &wg, name: "consumer -- 2"}
	err = consumer2.Consume(ctx, consunerHandler2.handle)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	wg.Add(2)
	err = testPublisher.Publish(
		ctx,
		msgbus.Metadata{Id: testId},
		[]byte("test payload 2222"),
	)
	require.NoError(t, err)

	err = testPublisher.Publish(
		ctx,
		msgbus.Metadata{},
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

func (h handlerSvc) handle(ctx context.Context, subject string, _ msgbus.Metadata, data []byte) error {
	defer h.wg.Done()
	// ctx, span := logger.Span(ctx, "ConsumerLoop.handler", nil)
	authInfo := roles.AuthInfoFromCtx(ctx)
	assert.Equalf(h.t, "mockuser", authInfo.Subject, "authInfo is not set %#v", authInfo)
	log.Log().InfoCtx(ctx, "asdasd", slog.String("group", h.name), slog.String("data", string(data)))
	// span.End(nil)
	return nil
}
