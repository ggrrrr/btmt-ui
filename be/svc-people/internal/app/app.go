package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/pin"
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
		PinParse(ctx context.Context, pin string) (*ddd.PinValidation, error)
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
		logger.Warn().Msg("use mock AppPolices")
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
	if err := a.appPolices.CanDo(authInfo.Tenant, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return err
	}
	if p.Id != "" {
		old, _ := a.repoPeople.GetById(ctx, p.Id)
		if old != nil {
			return app.BadRequestError("person with id exists", nil)
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
	parseEGN(p)

	err = a.repoPeople.Save(ctx, p)
	if err != nil {
		return err
	}
	logger.DebugCtx(ctx).Any("data", p).Msg("Save")
	return nil
}

func (a *application) GetById(ctx context.Context, id string) (*ddd.Person, error) {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(authInfo.Tenant, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return nil, err
	}
	logger.InfoCtx(ctx).Str("id", id).Msg("GetById")
	person, err := a.repoPeople.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if person == nil {
		return nil, app.ItemNotFoundError("person", id)
	}
	return person, nil
}

func (a *application) List(ctx context.Context, filters Filters) ([]ddd.Person, error) {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(authInfo.Tenant, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		logger.ErrorCtx(ctx, err).Msg("List")
		return nil, err
	}
	logger.DebugCtx(ctx).Any("filters", filters).Msg("List")

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
		logger.ErrorCtx(ctx, err).Msg("List")

		return nil, err
	}
	return out, nil
}

func (a *application) Update(ctx context.Context, p *ddd.Person) error {
	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(authInfo.Tenant, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return err
	}
	logger.DebugCtx(ctx).Any("person", p).Msg("Update")
	err := a.repoPeople.Update(ctx, p)
	return err
}

func (*application) PinParse(ctx context.Context, number string) (*ddd.PinValidation, error) {
	info, err := pin.Parse(number)
	if err != nil {
		return nil, err
	}
	return &ddd.PinValidation{
		Dob:    info.Dob,
		Gender: info.Gender,
	}, nil
}

func parseEGN(person *ddd.Person) {
	res, err := pin.Parse(person.IdNumbers["EGN"])
	if err != nil {
		logger.Error(err).Any("EGN", person)
		return
	}
	if person.DOB == nil {
		person.DOB = &ddd.Dob{}
	}
	if res.Dob.Year > 0 {
		person.DOB.Year = res.Dob.Year
	}
	if res.Dob.Month > 0 {
		person.DOB.Month = res.Dob.Month
	}
	if res.Dob.Day > 0 {
		person.DOB.Day = res.Dob.Day
	}
	if len(res.Gender) > 0 {
		person.Gender = res.Gender
	}
	logger.Debug().Any("pin", person)
}
