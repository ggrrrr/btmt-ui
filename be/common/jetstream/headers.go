package jetstream

import (
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
)

const (
	contentTypeKey   string = "content-type"
	replyCounterKey  string = "reply-counter"
	replyAfterSecKey string = "reply-after-sec"
	replyTopicKey    string = "reply-topic"
)

func createMessageMD(from jetstream.Msg) msgbus.Metadata {
	fromHeader := from.Headers()

	md := msgbus.Metadata{}
	replyTopic, _ := getStr(replyTopicKey, fromHeader)
	replyAfter, _ := getInt(replyAfterSecKey, fromHeader)
	replyCounter, ok := getInt(replyCounterKey, fromHeader)
	if ok {
		md.RetryCounter = replyCounter
		md.RetryAfter = time.Duration(replyAfter) * time.Second
		md.RetryTopic = replyTopic
	}

	contentType, _ := getStr(contentTypeKey, fromHeader)
	md.ContentType = contentType
	return md
}

func getInt(key string, h nats.Header) (int, bool) {
	s, ok := h[key]
	if !ok {
		return 0, false
	}
	i, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, false
	}
	return i, true
}

func getStr(key string, h nats.Header) (string, bool) {
	s, ok := h[key]
	if !ok {
		return "", false
	}
	return s[0], true
}

func injectHeaders(md msgbus.Metadata, natsMsg publishMsg) {
	if md.ContentType != "" {
		natsMsg.Set(contentTypeKey, md.ContentType)
	}
	if md.RetryCounter > 0 {
		natsMsg.Set(replyCounterKey, strconv.Itoa(md.RetryCounter))
	}
	if len(md.RetryTopic) > 0 {
		natsMsg.Set(replyTopicKey, md.RetryTopic)
	}
	if md.RetryAfter > 0 {
		natsMsg.Set(replyAfterSecKey, strconv.Itoa(int(md.RetryAfter.Seconds())))
	}
}
