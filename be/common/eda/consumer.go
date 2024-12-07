package eda

import "context"

type (
	Consumer interface {
		Consume(ctx context.Context, handler func(ctx context.Context, subject string, data []byte)) error
		Shutdown() error
	}
)
