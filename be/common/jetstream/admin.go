package jetstream

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

func (c *NatsConnection) CreateStream(ctx context.Context, name string, descr string, subjects []string) (jetstream.Stream, error) {

	stream, err := c.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        name,
		Description: descr,
		Subjects:    subjects,
		Metadata: map[string]string{
			"agent": "btmt",
		},
	})
	return stream, err
}
func (c *NatsConnection) PruneStream(ctx context.Context, name string) error {
	return c.js.DeleteStream(ctx, name)
}

func (c *NatsConnection) Stream(ctx context.Context, name string) (jetstream.Stream, error) {
	s, err := c.js.Stream(ctx, name)
	return s, err
}
