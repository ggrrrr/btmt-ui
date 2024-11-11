package repo

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

type (
	Repo struct {
		collection string
		db         mgo.Repo
	}
)

func (r *Repo) Save(ctx context.Context, template *ddd.Template) (err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.Save", nil, logger.TraceKVString("template.name", template.Name))
	defer func() {
		span.End(err)
	}()

	newTmpl, err := fromTemplate(template)
	if err != nil {
		return
	}
	_, err = r.db.InsertOne(ctx, r.collection, newTmpl)
	if err != nil {
		return
	}
	template.Id = newTmpl.Id.Hex()

	return nil
}

func (r *Repo) List(ctx context.Context, filter app.FilterFactory) (result []ddd.Template, err error) {
	_, span := logger.Span(ctx, "repo.List", nil)
	defer func() {
		span.End(err)
	}()
	logger.InfoCtx(ctx).Msg("repo.List")

	cur, err := r.db.Find(ctx, r.collection, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("collection.Find %w", err)
	}
	defer cur.Close(ctx)
	if cur.Err() != nil {
		return nil, fmt.Errorf("collection.cursor %w", err)
	}

	var out = make([]ddd.Template, 0)
	for cur.Next(ctx) {
		if cur.Err() != nil {
			return nil, fmt.Errorf("cursor.Next %w", err)
		}

		current := internalTmpl{}

		err = cur.Decode(&current)
		if err != nil {
			return out, fmt.Errorf("unable to decode data %w", err)
		}

		out = append(out, current.toTemplate())
	}
	return out, nil
}

func (r *Repo) GetById(ctx context.Context, fromId string) (*ddd.Template, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "repo.GetById", nil, logger.TraceKVString("id", fromId), logger.TraceKVString("collection", r.collection))
	defer func() {
		span.End(err)
	}()

	id, err := mgo.ConvertFromId(fromId)
	if err != nil {
		return nil, fmt.Errorf("bad id[%s]", fromId)
	}

	res := r.db.FindOne(ctx, r.collection, bson.M{"_id": id})
	logger.DebugCtx(ctx).
		Str("collection", r.collection).
		Str("fromId", fromId).
		Str("id.Hex", id.Hex()).
		Any("id", id).
		Send()

	if res.Err() != nil {
		fmt.Printf("err %+v \n", err)
		err = res.Err()
		return nil, err
	}

	internal := internalTmpl{}
	err = res.Decode(&internal)
	if err != nil {
		return nil, app.SystemError("unable to decode record", err)
	}

	out := internal.toTemplate()

	return &out, err
}

func (r *Repo) Update(ctx context.Context, template *ddd.Template) (err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.Save", nil,
		logger.TraceKVString("template.id", template.Id),
		logger.TraceKVString("template.name", template.Name),
	)
	defer func() {
		span.End(err)
	}()

	id, err := mgo.ConvertFromId(template.Id)
	if err != nil {
		return
	}

	setReq := bson.M{}
	if len(template.Name) > 0 {
		setReq["name"] = template.Name
	}
	if len(template.ContentType) > 0 {
		setReq["content_type"] = template.ContentType
	}
	if len(template.Body) > 0 {
		setReq["body"] = template.Body
	}
	if len(template.Labels) > 0 {
		setReq["labels"] = template.Labels
	}

	setReq["created_at"] = mgo.FromTimeOrNow(time.Now())

	updateReq := bson.M{
		"$set": setReq,
	}

	logger.DebugCtx(ctx).
		Any("updateReq", updateReq).
		Str("id", template.Id).
		Msg("UpdateByID")
	resp, err := r.db.UpdateByID(ctx, r.collection, id, updateReq)
	if err != nil {
		return
	}

	logger.InfoCtx(ctx).
		Any("id", template.Id).
		Any("matchedCount", resp.MatchedCount).
		Msg("update")

	return

}
