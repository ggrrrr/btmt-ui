package jetstream

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const authHeaderName string = "authorization"

type natsConn struct {
	conn *nats.Conn
	js   jetstream.JetStream
}

func connect(url string) (*natsConn, error) {
	cn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(cn,
		jetstream.WithPublishAsyncMaxPending(256),
	)
	if err != nil {
		return nil, err
	}

	return &natsConn{
		conn: cn,
		js:   js,
	}, nil

}

func (c *natsConn) shutdown() error {
	if c.conn == nil {
		return nil
	}
	if c.conn.IsClosed() {
		return nil
	}
	c.conn.Close()
	return nil

}
