package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/state"
)

type StateStore struct {
	tracer   tracer.OTelTracer
	natsConn *NatsConnection
	bucket   string
	kv       jetstream.KeyValue
}

func NewStateStore(ctx context.Context, cfg Config, objectType state.EntityType) (*StateStore, error) {
	conn, err := Connect(cfg)
	if err != nil {
		return nil, err
	}

	kv, err := conn.js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      objectType.String(),
		Description: "",
		History:     10,
	})
	if err != nil {
		return nil, fmt.Errorf("create store[%s], %w", objectType, err)
	}

	setter := &StateStore{
		tracer:   tracer.Tracer(otelScope),
		natsConn: conn,
		bucket:   objectType.String(),
		kv:       kv,
	}

	return setter, nil
}
