package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb"
)

type (
	internalTmpl struct {
		Id          primitive.ObjectID `bson:"_id"`
		Labels      []string           `bson:"labels"`
		ContentType string             `bson:"content_type"`
		Name        string             `bson:"name"`
		Images      []string           `bson:"images"`
		Files       map[string]string  `bson:"files"`
		Body        string             `bson:"body"`
		CreatedAt   primitive.DateTime `bson:"created_at"`
		UpdatedAt   primitive.DateTime `bson:"updated_at"`
	}
)

func (from internalTmpl) toTemplate() *tmplpb.Template {
	return &tmplpb.Template{
		Id:          from.Id.Hex(),
		ContentType: from.ContentType,
		Name:        from.Name,
		Labels:      from.Labels,
		Images:      from.Images,
		Files:       from.Files,
		Body:        from.Body,
		CreatedAt:   timestamppb.New(from.CreatedAt.Time()),
		UpdatedAt:   timestamppb.New(from.UpdatedAt.Time()),
	}
}

func fromTemplate(from *tmplpb.Template) (internalTmpl, error) {

	id, err := mgo.ConvertFromId(from.Id)
	if err != nil {
		return internalTmpl{}, err
	}

	return internalTmpl{
		Id:          id,
		ContentType: from.ContentType,
		Labels:      from.Labels,
		Name:        from.Name,
		Body:        from.Body,
		Images:      from.Images,
		Files:       from.Files,
		CreatedAt:   mgo.FromTimeOrNow(from.CreatedAt.AsTime()),
		UpdatedAt:   mgo.FromTimeOrNow(from.UpdatedAt.AsTime()),
	}, nil
}
