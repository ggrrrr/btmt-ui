package state

import (
	"context"
)

type (
	StateGetter interface {
		Fetch(ctx context.Context, key string) (EntityState, error)
		History(ctx context.Context, key string) ([]EntityState, error)
	}

	StatePusher interface {
		// Pushes new entity to store and return revision
		Push(ctx context.Context, entity NewEntity) (uint64, error)
	}

	StateStore interface {
		StateGetter
		StatePusher
	}
)
