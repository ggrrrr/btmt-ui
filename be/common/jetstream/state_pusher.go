package jetstream

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/state"
)

var _ (state.StatePusher) = (*StateStore)(nil)

func (s *StateStore) Push(ctx context.Context, object state.NewEntity) (rev uint64, err error) {

	ctx, span := s.tracer.SpanWithAttributes(ctx, "jetstream.store.Push",
		slog.String("store.entity.id", object.Key),
		slog.String("store.entity.type", s.bucket),
	)
	defer func() {
		span.End(err)
	}()

	rev, err = s.kv.Put(ctx, object.Key, object.Value)
	if err != nil {
		err = fmt.Errorf("kv[%s] %w", s.bucket, err)
		return
	}
	log.Log().InfoCtx(ctx, "jetstream.store.Push",
		slog.String("rev", strconv.FormatUint(rev, 10)))
	return
}
