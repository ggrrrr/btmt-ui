package jetstream

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const authHeaderName string = "authorization"

type (
	Config struct {
		URL string
	}

	NatsConnection struct {
		conn *nats.Conn
		js   jetstream.JetStream
	}
)

func Connect(cfg Config) (*NatsConnection, error) {
	cn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(cn,
		jetstream.WithPublishAsyncMaxPending(256),
	)
	if err != nil {
		return nil, err
	}

	return &NatsConnection{
		conn: cn,
		js:   js,
	}, nil

}

func (c *NatsConnection) shutdown() error {
	if c.conn == nil {
		return nil
	}
	if c.conn.IsClosed() {
		return nil
	}
	c.conn.Close()
	return nil

}
