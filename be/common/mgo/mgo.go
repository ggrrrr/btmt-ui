package mgo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type (
	Config struct {
		TTL        time.Duration
		Collection string
		User       string
		Password   string
		Database   string
		Uri        string
		Debug      string
	}

	Repo interface {
		Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
		FindOne(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult

		InsertOne(ctx context.Context, collection string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
		UpdateByID(ctx context.Context, collection string, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

		Collection(c string) *mongo.Collection
	}

	repo struct {
		ctx    context.Context
		client *mongo.Client
		db     *mongo.Database
		cfg    Config
	}
)

var _ (Repo) = (*repo)(nil)

func New(ctx context.Context, cfg Config) (*repo, error) {
	logger.Info().
		Str("user", cfg.User).
		Str("database", cfg.Database).
		Str("uri", cfg.Uri).
		Str("collection", cfg.Collection).
		Str("debug", cfg.Debug).
		Any("ttl", cfg.TTL.Seconds()).
		Msg("New")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username:      cfg.User,
		Password:      cfg.Password,
		AuthMechanism: "SCRAM-SHA-1",
		// AuthMechanismProperties: {},
		// AuthMechanism: "SCRAM-SHA-256",
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client()
	opts.ApplyURI(cfg.Uri).SetServerAPIOptions(serverAPI)
	opts.SetAuth(credential)

	if cfg.Debug == "console" {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				// logger.Log().Info().Any("command", evt.Command.String()).Msg("mongo.event")
				fmt.Printf("mongo.event: %+v \n", evt.Command)
			},
		}
		opts.SetMonitor(cmdMonitor)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.ErrorCtx(ctx, err).
			Str("user", cfg.User).
			Str("database", cfg.Database).
			Str("uri", cfg.Uri).
			Str("collection", cfg.Collection).
			Str("debug", cfg.Debug).
			Any("ttl", cfg.TTL.Seconds()).
			Msg("Error")

		return nil, err
	}
	logger.Info().Msg("Connected")

	// d, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	logger.Error(err).
	// 		Str("user", cfg.User).
	// 		Str("password", cfg.Password).
	// 		Str("database", cfg.Database).
	// 		Str("uri", cfg.Uri).
	// 		Str("collection", cfg.Collection).
	// 		Str("debug", cfg.Debug).
	// 		Any("ttl", cfg.TTL.Seconds()).
	// 		Msg("ListDatabaseNames")

	// 	return nil, err
	// }
	// fmt.Printf("ListDatabaseNames: %v, \n", d)

	db := client.Database(cfg.Database)
	// cc, err := db.ListCollectionNames(ctx, bson.M{})
	// if err != nil {
	// 	logger.Error(err).
	// 		Str("user", cfg.User).
	// 		Str("password", cfg.Password).
	// 		Str("database", cfg.Database).
	// 		Str("uri", cfg.Uri).
	// 		Str("collection", cfg.Collection).
	// 		Str("debug", cfg.Debug).
	// 		Any("ttl", cfg.TTL.Seconds()).
	// 		Msg("ListCollectionNames")

	// 	return nil, err
	// }
	// fmt.Printf("ListCollectionNames: %v, \n", cc)

	logger.Info().Msg("ok")
	return &repo{
		cfg:    cfg,
		client: client,
		ctx:    ctx,
		db:     db,
	}, nil
}

func (r *repo) Close(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.TTL)
	defer cancel()

	err := r.db.Client().Disconnect(ctx)
	if err != nil {
		logger.Error(err).Msg("db.Close")
	}
	logger.Info().Msg("db.Closed")
}

func (r *repo) Collection(c string) *mongo.Collection {
	return r.db.Collection(c)
}

func (r *repo) Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.TTL)
	defer cancel()
	col := r.db.Collection(collection)
	return col.Find(ctx, filter, opts...)
}

func (r *repo) InsertOne(ctx context.Context, collection string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.TTL)
	defer cancel()

	col := r.db.Collection(collection)
	return col.InsertOne(ctx, document, opts...)
}

func (r *repo) UpdateByID(ctx context.Context, collection string, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.TTL)
	defer cancel()

	col := r.db.Collection(r.cfg.Collection)
	return col.UpdateByID(ctx, id, update, opts...)
}

func (r *repo) FindOne(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.TTL)
	defer cancel()

	col := r.db.Collection(r.cfg.Collection)
	return col.FindOne(ctx, filter, opts...)
}

func (r *repo) DB() *mongo.Database {
	return r.db
}
