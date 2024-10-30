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
		Name        string             `bson:"Name"`
		Body        string             `bson:"body"`
		CreatedAt   primitive.DateTime `bson:"created_at"`
	}
)

func New(collection string, db mgo.Repo) *Repo {
	return &Repo{
		collection: collection,
		db:         db,
	}
}

func (from internalTmpl) toTemplate() ddd.Template {
	return ddd.Template{
		Id:          from.Id.Hex(),
		ContentType: from.ContentType,
		Name:        from.Name,
		Labels:      from.Labels,
		Body:        from.Body,
		CreatedAt:   from.CreatedAt.Time(),
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
		CreatedAt:   mgo.FromTimeOrNow(from.CreatedAt),
	}, nil
}
