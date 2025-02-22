package mgo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type (
	Config struct {
		TTL        time.Duration `env:"MGO_TTL" envDefault:"1s"`
		Collection string        `env:"MGO_COLLECTION"`
		User       string        `env:"MGO_USER"`
		Password   string        `env:"MGO_PASSWORD"`
		Database   string        `env:"MGO_DATABASE"`
		Uri        string        `env:"MGO_URI"`
		Host       string        `env:"MGO_HOST"`
		Debug      string        `env:"MGO_DEBUG"`
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
		ttl    time.Duration
		// cfg    Config
	}
)

var _ (Repo) = (*repo)(nil)

func New(ctx context.Context, cfg Config) (*repo, error) {
	if cfg.TTL == 0 {
		cfg.TTL = time.Second
	}

	logger.Info().
		Str("user", cfg.User).
		Str("database", cfg.Database).
		Str("host", cfg.Host).
		Str("collection", cfg.Collection).
		Str("debug", cfg.Debug).
		Any("ttl", cfg.TTL.Seconds()).
		Msg("New")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var uri = ""
	credential := options.Credential{}
	credential.Username = cfg.User
	credential.Password = cfg.Password

	if len(cfg.Host) > 0 {
		uri = fmt.Sprintf("mongodb://%s:%s@%s/?ssl=false&authSource=admin",
			cfg.User,
			cfg.Password,
			cfg.Host,
		)
	}
	if len(cfg.Uri) > 0 {
		uri = cfg.Uri
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(uri).SetServerAPIOptions(serverAPI)
	opts.SetTimeout(cfg.TTL)

	if len(credential.Username) > 0 {
		opts.SetAuth(credential)
	}

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

	db := client.Database(cfg.Database)

	logger.Info().Msg("ok")
	return &repo{
		// cfg:    cfg,
		client: client,
		ctx:    ctx,
		db:     db,
		ttl:    cfg.TTL,
	}, nil
}

func (r *repo) Close(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, r.ttl)
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
	ctx, cancel := context.WithTimeout(ctx, r.ttl)
	defer cancel()
	col := r.db.Collection(collection)
	return col.Find(ctx, filter, opts...)
}

func (r *repo) InsertOne(ctx context.Context, collection string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ttl)
	defer cancel()

	col := r.db.Collection(collection)
	return col.InsertOne(ctx, document, opts...)
}

func (r *repo) UpdateByID(ctx context.Context, collection string, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ttl)
	defer cancel()

	col := r.db.Collection(collection)
	return col.UpdateByID(ctx, id, update, opts...)
}

func (r *repo) FindOne(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, r.ttl)
	defer cancel()

	col := r.db.Collection(collection)
	return col.FindOne(ctx, filter, opts...)
}

func (r *repo) DB() *mongo.Database {
	return r.db
}
