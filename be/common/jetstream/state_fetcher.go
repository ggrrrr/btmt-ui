package jetstream

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ggrrrr/btmt-ui/be/common/state"
)

var _ (state.StateFetcher) = (*StateStore)(nil)

func (s *StateStore) History(ctx context.Context, key string) (states []state.EntityState, err error) {
	ctx, span := s.tracer.SpanWithAttributes(ctx, "jetstream.store.History",
		slog.String("store.entity.id", key),
		slog.String("store.entity.type", s.bucket),
	)
	defer func() {
		span.End(err)
	}()

	entrys, err := s.kv.History(ctx, key)
	if err != nil {
		err = fmt.Errorf("History[%s].Get(%s) %w", s.bucket, key, err)
		return
	}

	states = make([]state.EntityState, 0, len(entrys))

	for k := range entrys {
		entry := entrys[k]

		entity := state.EntityState{
			Revision:  entry.Revision(),
			Key:       key,
			Value:     entry.Value(),
			CreatedAt: entry.Created(),
		}
		states = append(states, entity)
	}

	return

}

func (s *StateStore) Fetch(ctx context.Context, key string) (entity state.EntityState, err error) {
	ctx, span := s.tracer.SpanWithAttributes(ctx, "jetstream.store.Fetch",
		slog.String("entity.id", key),
		slog.String("entity.type", s.bucket),
	)
	defer func() {
		span.End(err)
	}()

	entry, err := s.kv.Get(ctx, key)
	if err != nil {
		err = fmt.Errorf("StateGetter[%s].Fetch(%s) %w", s.bucket, key, err)
		return
	}

	if entry == nil {
		err = fmt.Errorf("not found")
		return
	}

	entity = state.EntityState{
		Revision:  entry.Revision(),
		Key:       key,
		Value:     entry.Value(),
		CreatedAt: entry.Created(),
	}
	return
}
