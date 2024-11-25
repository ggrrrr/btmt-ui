package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
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
		BlobId      string             `bson:"blob_id"`
		CreatedAt   primitive.DateTime `bson:"created_at"`
		UpdatedAt   primitive.DateTime `bson:"updated_at"`
	}
)

func (from internalTmpl) toTemplate() ddd.Template {
	return ddd.Template{
		Id:          from.Id.Hex(),
		ContentType: from.ContentType,
		Name:        from.Name,
		Labels:      from.Labels,
		Images:      from.Images,
		Files:       from.Files,
		Body:        from.Body,
		BlobId:      from.BlobId,
		CreatedAt:   from.CreatedAt.Time(),
		UpdatedAt:   from.UpdatedAt.Time(),
	}
}

func fromTemplate(from *ddd.Template) (internalTmpl, error) {

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
		BlobId:      from.BlobId,
		CreatedAt:   mgo.FromTimeOrNow(from.CreatedAt),
		UpdatedAt:   mgo.FromTimeOrNow(from.UpdatedAt),
	}, nil
}
