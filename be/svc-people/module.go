package people

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/rest"
	peoplepbv1 "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

type (
	moduleCfg struct {
		MGO    mgo.Config `envPrefix:"PEOPLE_"`
		Broker jetstream.Config
	}
	Module struct {
		// service system.Service
		cfg moduleCfg
		app app.App
	}
)

var _ (system.Module) = (*Module)(nil)

func (*Module) Name() string {
	return "svc-people"
}

func (m *Module) Configure(ctx context.Context, s system.Service) (err error) {
	cfg := moduleCfg{}
	config.MustParse(&cfg)
	m.cfg = cfg

	db, err := mgo.Connect(ctx, m.cfg.MGO)
	if err != nil {
		log.Log().Error(err, "mgo")
		return err
	}

	natsConn, err := jetstream.Connect(m.cfg.Broker)

	s.Waiter().AddCleanup(func() {
		db.Close(context.Background())
	})

	stateStore, err := jetstream.NewStateStore(ctx, natsConn, state.EntityTypeFromProto(&peoplepbv1.Person{}))
	if err != nil {
		return err
	}

	m.app, err = app.New(
		app.WithStateStore(stateStore),
		app.WithPeopleRepo(repo.New(m.cfg.MGO.Collection, db)),
	)
	if err != nil {
		return err
	}

	s.MountHandler("/v1/people", rest.Router(rest.Handler(m.app)))

	return nil
}

func (*Module) Startup(ctx context.Context) (err error) {
	// return Root(ctx, s)
	return nil
}
