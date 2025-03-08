package mgo

import (
	"context"
	"time"
)

func ConnectForTest(collection string) (*repo, error) {
	cfg := Config{
		TTL:        10 * time.Second,
		Collection: collection,
		Database:   "tests",
		Uri:        "mongodb://admin:pass@localhost:27017/?retryWrites=true&w=majority&?authSource=admin",
	}
	cfg.Collection = collection

	return Connect(context.Background(), cfg)
}
