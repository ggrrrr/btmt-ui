package jetstream

import (
	"context"
	"log/slog"
	"strings"

	"github.com/nats-io/nats.go/jetstream"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
)

func (c *NatsConnection) CreateStream(ctx context.Context, name string, descr string, subjects []string) (jetstream.Stream, error) {
	log.Log().Info("CreateStream", slog.String("subjects", strings.Join(subjects, ",")))

	stream, err := c.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        name,
		Description: descr,
		Subjects:    subjects,
		Storage:     jetstream.FileStorage,
		Metadata: map[string]string{
			"agent": "btmt",
		},
		// Retention:            0,
		// MaxConsumers:         0,
		// MaxMsgs:              0,
		// MaxBytes:             0,
		// Discard:              0,
		// DiscardNewPerSubject: false,
		// MaxAge:               0,
		// MaxMsgsPerSubject:    0,
		// MaxMsgSize:           0,
		// Storage:              0,
		// Replicas:             0,
		// NoAck:                false,
		// Duplicates:           0,
		// Placement:            &jetstream.Placement{},
		// Mirror:               &jetstream.StreamSource{},
		// Sources:              []*jetstream.StreamSource{},
		// Sealed:               false,
		// DenyDelete:           false,
		// DenyPurge:            false,
		// AllowRollup:          false,
		// Compression:          jetstream.NoCompression,
		// FirstSeq:             0,
		// SubjectTransform:     &jetstream.SubjectTransformConfig{},
		// RePublish:            &jetstream.RePublish{},
		// AllowDirect:          false,
		// MirrorDirect:         false,
		// ConsumerLimits:       jetstream.StreamConsumerLimits{},

	})
	return stream, err
}
func (c *NatsConnection) PruneStream(ctx context.Context, name string) error {
	return c.js.DeleteStream(ctx, name)
}

func (c *NatsConnection) Stream(ctx context.Context, name string) (jetstream.Stream, error) {
	log.Log().Info("Stream", slog.String("name", name))

	s, err := c.js.Stream(ctx, name)
	return s, err
}
