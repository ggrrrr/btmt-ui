package tmpl

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob/awss3"
	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/repo"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/rest"
	tmplpbv1 "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

type Module struct{}

var _ (system.Module) = (*Module)(nil)

func (*Module) Name() string {
	return "svc-tmpl"
}

func (*Module) Startup(ctx context.Context, s system.Service) (err error) {
	return Root(ctx, s)
}

func Root(ctx context.Context, s system.Service) error {
	logger.Info().Msg("Root")

	db, err := mgo.New(ctx, s.Config().Mgo)
	if err != nil {
		logger.Error(err).Msg("db")
		return err
	}
	fn := func() {
		db.Close(ctx)
	}
	s.Waiter().Cleanup(fn)

	appRepo := repo.New("tmpl", db)

	blobClient, err := awss3.NewClient("test-bucket-1", awsclient.AwsConfig{
		Region:   "us-east-1",
		Endpoint: "http://localhost:4566",
	})
	if err != nil {
		return err
	}

	stateStore, err := jetstream.NewStateStore(ctx, jetstream.Config{
		URL: "localhost:4222",
	}, state.EntityTypeFromProto(&tmplpbv1.TemplateData{}))
	if err != nil {
		return err
	}

	a, err := app.New(
		app.WithBlobStore(blobClient),
		app.WithTmplRepo(appRepo),
		app.WithStateStore(stateStore),
	)
	if err != nil {
		return err
	}

	if s.Mux() == nil {
		return fmt.Errorf("system.Mux is nil")
	}

	restApp := rest.New(a)
	s.Mux().Mount("/tmpl", restApp.Router())

	return nil

}
