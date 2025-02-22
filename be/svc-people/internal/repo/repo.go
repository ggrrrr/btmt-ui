package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/mgo.v2/bson"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

type (
	repo struct {
		collection string
		db         mgo.Repo
	}

	// Repo interface {
	// 	Close() error
	// 	Save(ctx context.Context, p *ddd.Person) error
	// }
)

var _ (ddd.PeopleRepo) = (*repo)(nil)

func New(collection string, db mgo.Repo) *repo {
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
	if _, err := c.Indexes().CreateOne(ctx, mod); err != nil {
		logger.ErrorCtx(ctx, err).Msg("CreateIndex")
	}
}

func (r *repo) Save(ctx context.Context, p *peoplepbv1.Person) (err error) {
	ctx, span := logger.Span(ctx, "repo.Save", p)
	defer func() {
		span.End(err)
	}()

	p.CreatedAt = timestamppb.Now()
	newPerson, err := fromPerson(p)
	if err != nil {
		return
	}
	_, err = r.db.InsertOne(ctx, r.collection, newPerson)
	if err != nil {
		return
	}
	p.Id = newPerson.Id.Hex()

	return nil
}

func (r *repo) Update(ctx context.Context, p *peoplepbv1.Person) (err error) {
	_, span := logger.Span(ctx, "repo.Update", p)
	defer func() {
		span.End(err)
	}()

	if len(p.Id) == 0 {
		err = app.BadRequestError("person.id is empty", nil)
		return
	}

	p.UpdatedAt = timestamppb.Now()

	newPerson, err := fromPerson(p)
	if err != nil {
		return err
	}
	setReq := bson.M{}

	if len(newPerson.IdNumbers) > 0 {
		setReq[FieldIdNumbers] = newPerson.IdNumbers
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

	setReq[FieldUpdatedAt] = newPerson.UpdatedAt

	if len(setReq) == 0 {
		err = app.BadRequestError("empty person update", nil)
		return
	}
	updateReq := bson.M{
		"$set": setReq,
	}
	logger.DebugCtx(ctx).Any("updateReq", updateReq).Msg("UpdateByID")
	resp, err := r.db.UpdateByID(ctx, r.collection, newPerson.Id, updateReq)
	if err != nil {
		return
	}

	logger.InfoCtx(ctx).
		Any("id", newPerson.Id).
		Any("matchedCount", resp.MatchedCount).
		Msg("update")

	return nil
}

func (r *repo) List(ctx context.Context, filter app.FilterFactory) (result []*peoplepbv1.Person, err error) {
	_, span := logger.Span(ctx, "repo.List", nil)
	defer func() {
		span.End(err)
	}()

	if filter == nil {
		result, err = r.list(ctx, bson.M{})
		return
	}
	result, err = r.list(ctx, filter.Create())
	return
}

func (r *repo) GetById(ctx context.Context, fromId string) (result *peoplepbv1.Person, err error) {
	_, span := logger.SpanWithAttributes(ctx, "repo.GetById", nil, logger.TraceKVString("id", fromId))
	defer func() {
		span.End(err)
	}()

	id, err := mgo.ConvertFromId(fromId)
	if err != nil {
		return
	}

	logger.DebugCtx(ctx).
		Str("fromId", fromId).
		Str("id.Hex", id.Hex()).
		Str("id", id.String()).
		Send()

	var out person
	res := r.db.FindOne(ctx, r.collection, bson.M{"_id": id})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, nil
		}
		err = res.Err()
		return
	}
	err = res.Decode(&out)
	if err != nil {
		return nil, app.SystemError("unable to decode record", err)
	}
	result = out.toProto()
	return
}

func (r *repo) list(ctx context.Context, filter any) ([]*peoplepbv1.Person, error) {

	cur, err := r.db.Find(ctx, r.collection, filter)
	if err != nil {
		return nil, fmt.Errorf("collect.Find %w", err)
	}
	defer cur.Close(ctx)
	if cur.Err() != nil {
		return nil, fmt.Errorf("collect.Find.cursor %w", err)
	}

	var out = make([]*peoplepbv1.Person, 0)
	for cur.Next(context.Background()) {
		if cur.Err() != nil {
			return nil, fmt.Errorf("collect.Find.cursor.Next %w", err)
		}
		var result *person
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("collect.Find decode %w", err)
		}
		out = append(out, result.toProto())
	}

	return out, nil
}
