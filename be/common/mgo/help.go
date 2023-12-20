package mgo

import (
	"os"
	"time"
)

func MgoTestCfg() Config {
	cfg := Config{
		TTL:        10 * time.Second,
		Collection: os.Getenv("MGO_COLLECTION"),
		User:       os.Getenv("MGO_USERNAME"),
		Password:   os.Getenv("MGO_PASSWORD"),
		Database:   os.Getenv("MGO_DATABASE"),
		Uri:        os.Getenv("MGO_URI"),
		Host:       os.Getenv("MGO_HOST"),
	}
	return cfg
}
