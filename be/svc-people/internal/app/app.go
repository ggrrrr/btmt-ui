package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	"github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb"
)

type (
	AppConfiguration func(a *application) error

	Filters map[string][]string

	application struct {
		repoPeople ddd.PeopleRepo
		appPolices roles.AppPolices
	}

	App interface {
		Save(ctx context.Context, p *ddd.Person) error
		List(ctx context.Context, filters Filters) ([]ddd.Person, error)
		GetById(ctx context.Context, id string) (*ddd.Person, error)
		Update(ctx context.Context, p *ddd.Person) error
	}
)

var (
	FilterTexts  string = "texts"
	FilterPhones string = "phones"
	FilterPINs   string = "pins"
	FilterLabels string = "labels"
	FilterAttrs  string = "attrs"
)

var _ (App) = (*application)(nil)

func New(cfgs ...AppConfiguration) (*application, error) {
	a := &application{}
	for _, c := range cfgs {
		err := c(a)
		if err != nil {
			return nil, err
		}
	}
	if a.appPolices == nil {
		logger.Log().Warn().Msg("use mock AppPolices")
		a.appPolices = roles.NewAppPolices()
	}
	return a, nil
}

func WithPeopleRepo(repo ddd.PeopleRepo) AppConfiguration {
	return func(a *application) error {
		a.repoPeople = repo
		return nil
	}
}

func WithAppPolicies(appPolices roles.AppPolices) AppConfiguration {
	return func(a *application) error {
		a.appPolices = appPolices
		return nil
	}
}

func (a *application) Save(ctx context.Context, p *ddd.Person) error {
	var err error
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return err
	}
	if p.Id != "" {
		old, _ := a.repoPeople.GetById(ctx, p.Id)
		if old != nil {
			return app.ErrorBadRequest("person with id exists", nil)
		}
	}
	if p.Phones == nil {
		p.Phones = map[string]string{}
	}
	if p.Labels == nil {
		p.Labels = []string{}
	}
	if p.Attr == nil {
		p.Attr = map[string]string{}
	}
	err = a.repoPeople.Save(ctx, p)
	if err != nil {
		return err
	}
	logger.Log().Debug().Any("data", p).Any("trace", logger.LogTraceData(ctx)).Msg("Save")
	return nil
}

func (a *application) GetById(ctx context.Context, id string) (*ddd.Person, error) {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return nil, err
	}
	logger.Log().Info().Str("id", id).Any("trace", logger.LogTraceData(ctx)).Msg("GetById")
	return a.repoPeople.GetById(ctx, id)
}

func (a *application) List(ctx context.Context, filters Filters) ([]ddd.Person, error) {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		logger.Log().Error().Err(err).Any("filters", filters).Any("trace", logger.LogTraceData(ctx)).Msg("List")
		return nil, err
	}
	logger.Log().Debug().Any("filters", filters).Any("trace", logger.LogTraceData(ctx)).Msg("List")

	filtersFunc := []repo.AddFilterFunc{}
	if len(filters[FilterPINs]) > 0 {
		filtersFunc = append(filtersFunc, repo.AddPINs(filters[FilterPINs]...))
	}
	if len(filters[FilterTexts]) > 0 {
		filtersFunc = append(filtersFunc, repo.AddTexts(filters[FilterTexts]...))
	}
	if len(filters[FilterLabels]) > 0 {
		filtersFunc = append(filtersFunc, repo.AddLabels(filters[FilterLabels]...))
	}
	if len(filters[FilterPhones]) > 0 {
		filtersFunc = append(filtersFunc, repo.AddPhones(filters[FilterPhones]...))
	}

	ff, err := repo.NewFilter(filtersFunc...)
	if err != nil {
		return nil, err
	}

	out, err := a.repoPeople.List(ctx, ff)
	if err != nil {
		logger.Log().Error().Err(err).Any("trace", logger.LogTraceData(ctx)).Msg("List")

		return nil, err
	}
	return out, nil
}

func (a *application) Update(ctx context.Context, p *ddd.Person) error {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return err
	}
	logger.Log().Debug().Any("trace", logger.LogTraceData(ctx)).Msg("Update")
	err := a.repoPeople.Update(ctx, p)
	return err
}
