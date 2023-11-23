package repo

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/stackus/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/mongodb"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
)

type (
	repo struct {
		collection string
		db         mongodb.Repo
	}

	Repo interface {
		Close() error
		Save(ctx context.Context, p *ddd.Person) error
	}
)

var _ (ddd.PeopleRepo) = (*repo)(nil)

func New(collection string, db mongodb.Repo) *repo {
	return &repo{
		collection: collection,
		db:         db,
	}
}

func (r *repo) CreateIndex(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c := r.db.Collection(r.collection)
	mod := mongo.IndexModel{
		Keys: bson.M{
			"full_name": "text",
		},
		Options: nil,
	}
	c.Indexes().CreateOne(ctx, mod)
}

func (r *repo) Save(ctx context.Context, p *ddd.Person) error {
	p.CreatedTime = time.Now()
	newPerson, err := fromPerson(p)
	if err != nil {
		return err
	}
	_, err = r.db.InsertOne(ctx, r.collection, newPerson)
	if err != nil {
		return err
	}
	p.Id = newPerson.Id.Hex()

	return nil
}

func (r *repo) Update(ctx context.Context, p *ddd.Person) error {
	if len(p.Id) == 0 {
		return app.ErrorBadRequest("person.id is empty", nil)
	}
	newPerson, err := fromPerson(p)
	if err != nil {
		return err
	}
	setReq := bson.M{}

	if len(newPerson.PIN) > 0 {
		setReq[FieldPIN] = strings.Trim(newPerson.PIN, " ")
	}
	if len(newPerson.LoginEmail) > 0 {
		setReq[FieldLoginEmail] = strings.Trim(newPerson.LoginEmail, " ")
	}
	if len(newPerson.Emails) > 0 {
		setReq[FieldEmails] = newPerson.Emails
	}
	if len(newPerson.Name) > 0 {
		setReq[FieldName] = strings.Trim(newPerson.Name, " ")
	}
	if len(newPerson.FullName) > 0 {
		setReq[FieldFullName] = strings.Trim(newPerson.FullName, " ")
	}
	if newPerson.DOB != nil {
		if !newPerson.DOB.isZero() {
			setReq[FieldDoB] = newPerson.DOB
		}
	}
	if len(newPerson.Gender) > 0 {
		setReq[FieldGender] = strings.Trim(newPerson.Gender, " ")
	}
	if len(newPerson.Phones) > 0 {
		setReq[FieldPhones] = newPerson.Phones
	}
	if len(newPerson.Labels) > 0 {
		setReq[FieldLabels] = newPerson.Labels
	}
	if len(newPerson.Attr) > 0 {
		setReq[FieldAttr] = newPerson.Attr
	}

	if len(setReq) == 0 {
		return app.ErrorBadRequest("empty person", nil)
	}
	updateReq := bson.M{
		"$set": setReq,
	}
	logger.DebugCtx(ctx).Any("updateReq", updateReq).Msg("UpdateByID")
	resp, err := r.db.UpdateByID(ctx, r.collection, newPerson.Id, updateReq)
	if err != nil {
		return err
	}

	logger.InfoCtx(ctx).
		Any("id", newPerson.Id).
		Any("matchedCount", resp.MatchedCount).
		Msg("update")

	return nil
}

func (r *repo) List(ctx context.Context, filter ddd.FilterFactory) ([]ddd.Person, error) {
	if filter == nil {
		return r.list(ctx, bson.M{})
	}
	return r.list(ctx, filter.Create())
}

func (r *repo) GetById(ctx context.Context, fromId string) (*ddd.Person, error) {
	id, err := convertPersonId(fromId)
	if err != nil {
		return nil, err
	}

	var out person
	res := r.db.FindOne(ctx, r.collection, bson.M{"_id": id})

	if res.Err() != nil {
		if errors.As(res.Err(), &mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, res.Err()
	}
	err = res.Decode(&out)
	if err != nil {
		return nil, app.ErrorSystem("unable to decode record", err)
	}
	d := out.toPerson()
	return &d, nil
}

func (r *repo) list(ctx context.Context, filter any) ([]ddd.Person, error) {
	cur, err := r.db.Find(ctx, r.collection, filter)
	if err != nil {
		return nil, errors.Wrap(err, "collection.Find")
	}
	defer cur.Close(ctx)
	if cur.Err() != nil {
		return nil, errors.Wrap(err, "cursor")
	}

	var out = make([]ddd.Person, 0)
	for cur.Next(context.Background()) {
		if cur.Err() != nil {
			return nil, errors.Wrap(err, "cursor.Error")
		}
		var result person
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, errors.Wrap(err, "Decode")
		}
		out = append(out, result.toPerson())
	}

	return out, nil
}
