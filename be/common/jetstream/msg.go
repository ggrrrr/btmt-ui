package jetstream

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/propagation"

	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
)

type (
	publishMsg struct {
		md      msgbus.Metadata
		payload []byte
		msg     nats.Msg
	}

	jsMsg struct {
		msg jetstream.Msg
	}
)

var _ (propagation.TextMapCarrier) = (*publishMsg)(nil)

var _ (propagation.TextMapCarrier) = (*jsMsg)(nil)

// Get implements propagation.TextMapCarrier.
func (j jsMsg) Get(key string) string {
	if j.msg == nil {
		return ""
	}
	if j.msg.Headers() == nil {
		return ""
	}
	return j.msg.Headers().Get(key)
}

// Keys implements propagation.TextMapCarrier.
func (j jsMsg) Keys() []string {
	if j.msg == nil {
		return []string{}
	}
	if j.msg.Headers() == nil {
		return []string{}
	}
	out := make([]string, 0, len(j.msg.Headers()))
	for i := range j.msg.Headers() {
		out = append(out, i)
	}
	return out

}

// Set implements propagation.TextMapCarrier.
func (j jsMsg) Set(key string, value string) {
	if j.msg == nil {
		return
	}
	j.msg.Headers().Add(key, value)
}

// Get implements propagation.TextMapCarrier.
func (n publishMsg) Get(key string) string {
	if n.msg.Header == nil {
		return ""
	}
	return n.msg.Header.Get(key)
}

// Keys implements propagation.TextMapCarrier.
func (n *publishMsg) Keys() []string {
	if n.msg.Header == nil {
		return []string{}
	}
	out := make([]string, 0, len(n.msg.Header))
	for i := range n.msg.Header {
		out = append(out, i)
	}
	return out
}

// Set implements propagation.TextMapCarrier.
func (n *publishMsg) Set(key string, value string) {
	if n.msg.Header == nil {
		n.msg.Header = nats.Header{}
	}
	n.msg.Header.Add(key, value)
}
