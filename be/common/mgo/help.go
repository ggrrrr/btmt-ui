package mgo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
