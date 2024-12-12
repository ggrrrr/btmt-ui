package msgbus

import (
	"time"

	"github.com/google/uuid"
)

type (
	Metadata struct {
		Id           uuid.UUID
		ContentType  string
		UniqId       string
		OrderId      string
		RetryCounter int
		RetryAfter   time.Duration
		RetryTopic   string
	}
)
