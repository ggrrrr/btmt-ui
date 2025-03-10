package jetstream

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.be.jetstream"

const authHeaderName string = "authorization"

type (
	ConnOptionFunc func(a *NatsConnection) error

	Config struct {
		URL string `env:"NATS_URL"`
	}

	NatsConnection struct {
		tracer   tracer.OTelTracer
		conn     *nats.Conn
		js       jetstream.JetStream
		verifier token.Verifier
	}

	Connector interface {
	}
)

func Connect(cfg Config, opts ...ConnOptionFunc) (*NatsConnection, error) {
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

	n := &NatsConnection{
		tracer: tracer.Tracer(otelScope),
		conn:   cn,
		js:     js,
	}

	for _, f := range opts {
		err = f(n)
		if err != nil {
			return nil, err
		}
	}

	return n, nil

}

func (c *NatsConnection) Shutdown() error {
	if c.conn == nil {
		return nil
	}
	if c.conn.IsClosed() {
		return nil
	}
	c.conn.Close()
	return nil

}

func WithVerifier(verifier token.Verifier) ConnOptionFunc {
	return func(n *NatsConnection) error {
		n.verifier = verifier
		return nil
	}
}
