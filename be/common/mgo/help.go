package mgo

import (
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MgoTestCfg(collection string) Config {
	cfg := Config{
		TTL:        10 * time.Second,
		Collection: os.Getenv("MGO_COLLECTION"),
		User:       os.Getenv("MGO_USERNAME"),
		Password:   os.Getenv("MGO_PASSWORD"),
		Database:   os.Getenv("MGO_DATABASE"),
		Uri:        os.Getenv("MGO_URI"),
		Host:       os.Getenv("MGO_HOST"),
	}
	cfg.Collection = collection
	return cfg
}

func ConvertFromId(fromId string) (primitive.ObjectID, error) {
	if fromId == "" {
		return primitive.NewObjectID(), nil
	}
	return primitive.ObjectIDFromHex(fromId)
}

func FromTimeOrNow(fromTime time.Time) primitive.DateTime {
	if fromTime.IsZero() {
		return primitive.NewDateTimeFromTime(time.Now())
	}
	return primitive.NewDateTimeFromTime(fromTime)
}
