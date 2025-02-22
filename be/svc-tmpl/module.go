package tmpl

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob/awss3"
	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/repo"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/rest"
	tmplpbv1 "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

type (
	Config struct {
		MGO       mgo.Config `prefix:"TMPL_"`
		Broker    jetstream.Config
		BlobStore struct {
			BucketName string `env:"TMPL_STATE_BUCKET_NAME"`
			awsclient.Config
		}
	}

	Module struct {
		app    app.App
		cfg    Config
		system system.Service
	}
)

var _ (system.Module) = (*Module)(nil)

func (*Module) Name() string {
	return "svc-tmpl"
}
func (m *Module) Configure(ctx context.Context, s system.Service) (err error) {
	config.MustParse(&m.cfg)
	m.system = s
	db, err := mgo.New(ctx, m.cfg.MGO)
	if err != nil {
		return err
	}

	s.Waiter().AddCleanup(func() {
		db.Close(context.Background())
	})

	stateStore, err := jetstream.NewStateStore(ctx, m.cfg.Broker, state.EntityTypeFromProto(&tmplpbv1.TemplateData{}))
	if err != nil {
		return err
	}

	appRepo := repo.New("tmpl", db)

	blobClient, err := awss3.NewClient("test-bucket-1", awsclient.Config{
		Region:      "us-east-1",
		EndpointURL: "http://localhost:4566",
	})
	if err != nil {
		return err
	}

	m.app, err = app.New(
		app.WithBlobStore(blobClient),
		app.WithTmplRepo(appRepo),
		app.WithStateStore(stateStore),
	)
	if err != nil {
		return err
	}
	m.system.MountHandler("/v1/tmpl", rest.Router(rest.Handler(m.app)))
	return nil
}

func (m *Module) Startup(ctx context.Context) (err error) {

	return nil
}
