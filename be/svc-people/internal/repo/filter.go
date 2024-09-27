package repo

import (
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AddFilterFunc func(f *filter) error

	filter struct {
		pins   []string
		labels []string
		phones []string
		texts  []string
		// attr   []string
	}
)

var _ (app.FilterFactory) = (*filter)(nil)

func NewFilter(cfgs ...AddFilterFunc) (*filter, error) {
	out := &filter{}
	for _, c := range cfgs {
		err := c(out)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func AddTexts(texts ...string) AddFilterFunc {
	return func(f *filter) error {
		f.texts = append(f.texts, texts...)
		return nil
	}
}

func AddLabels(labels ...string) AddFilterFunc {
	return func(f *filter) error {
		f.labels = append(f.labels, labels...)
		return nil
	}
}

func AddPhones(phones ...string) AddFilterFunc {
	return func(f *filter) error {
		f.phones = append(f.phones, phones...)
		return nil
	}
}

func AddPINs(pins ...string) AddFilterFunc {
	return func(f *filter) error {
		f.pins = append(f.pins, pins...)
		return nil
	}
}

func (f *filter) Create() any {
	filter := bson.A{}
	if len(f.pins) > 0 {
		filter = append(filter, f.allPins())
	}
	if len(f.labels) > 0 {
		filter = append(filter, f.allLabels())
	}
	if len(f.phones) > 0 {
		filter = append(filter, f.allPhones())
	}
	if len(f.texts) > 0 {
		filter = append(filter, f.allTexts())
	}
	logger.Debug().Any("filters", filter).Msg("Create")
	if len(filter) > 0 {
		filter1 := bson.M{
			"$and": filter,
		}
		return filter1
	}
	return bson.M{}
}

func (f *filter) allTexts() primitive.M {
	fields := bson.A{}
	for _, txt := range f.texts {
		fields = append(fields, bson.D{
			{
				Key:   FieldName,
				Value: primitive.Regex{Pattern: txt, Options: "i"},
			},
		})
		fields = append(fields, bson.D{
			{
				Key:   FieldFullName,
				Value: primitive.Regex{Pattern: txt, Options: "i"},
			},
		})
		fields = append(fields, bson.D{
			{
				Key:   FieldEmails,
				Value: primitive.Regex{Pattern: txt, Options: "i"},
			},
		})
	}

	filter := primitive.M{
		"$or": fields,
	}

	return filter
}

func (f *filter) allPins() primitive.M {
	fields := bson.A{}
	for _, pin := range f.pins {
		fields = append(fields, bson.D{
			{
				Key:   FieldIdNumbers,
				Value: primitive.Regex{Pattern: fmt.Sprintf("%v", pin), Options: ""},
			},
		})
	}

	filter := primitive.M{
		"$or": fields,
	}

	return filter
}

func (f *filter) allLabels() primitive.M {
	fields := bson.A{}
	for _, label := range f.labels {
		fields = append(fields, bson.D{
			{
				Key:   FieldLabels,
				Value: primitive.Regex{Pattern: fmt.Sprintf("^%v", label), Options: ""},
			},
		})
	}

	filter := primitive.M{
		"$or": fields,
	}

	return filter
}

func (f *filter) allPhones() primitive.M {
	fields := bson.A{}
	for _, p := range f.phones {
		fields = append(fields, bson.D{
			{
				Key:   "phones",
				Value: primitive.Regex{Pattern: fmt.Sprintf("%v", p), Options: ""},
			},
		})
	}

	filter := primitive.M{
		"$or": fields,
	}

	return filter
}
