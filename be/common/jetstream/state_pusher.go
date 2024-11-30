package jetstream

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/state"
)

var _ (state.StatePusher) = (*StateStore)(nil)

func (s *StateStore) Push(ctx context.Context, object state.NewEntity) (rev uint64, err error) {

	ctx, span := logger.SpanWithAttributes(ctx, "jetstream.store.SetProto", nil,
		logger.TraceKVString("entity.id", object.Key),
		logger.TraceKVString("entity.type", s.bucket),
	)
	defer func() {
		span.End(err)
	}()

	rev, err = s.kv.Put(ctx, object.Key, object.Value)
	if err != nil {
		err = fmt.Errorf("kv[%s] %w", s.bucket, err)
		return
	}
	logger.InfoCtx(ctx).Str("rev", strconv.FormatUint(rev, 10)).Msg("jetstream.store.SetProto")
	return
}
