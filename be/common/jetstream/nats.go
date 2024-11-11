package jetstream

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type MsgHandler func(ctx context.Context, payload any) error

type Nats struct {
	conn *nats.Conn
	js   jetstream.JetStream
}

func Connect(url string) (*Nats, error) {
	cn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(cn, jetstream.WithPublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	return &Nats{
		conn: cn,
		js:   js,
	}, nil

}
