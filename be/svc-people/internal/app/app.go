package app

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/pin"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	peoplepb "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

type (
	AppConfiguration func(a *App) error

	Filters map[string][]string

	App struct {
		repoPeople ddd.PeopleRepo
		appPolices roles.AppPolices
		stateStore state.StateStore
	}

	// App interface {
	// 	Save(ctx context.Context, p *peoplepb.Person) error
	// 	List(ctx context.Context, filters Filters) ([]*peoplepb.Person, error)
	// 	GetById(ctx context.Context, id string) (*peoplepb.Person, error)
	// 	Update(ctx context.Context, p *peoplepb.Person) error
	// 	// PinParse(ctx context.Context, pin string) (*peoplepb.PinValidation, error)
	// }
)

var (
	FilterTexts  string = "texts"
	FilterPhones string = "phones"
	FilterPINs   string = "pins"
	FilterLabels string = "labels"
	FilterAttrs  string = "attrs"
)

// var _ (App) = (*App)(nil)

func New(cfgs ...AppConfiguration) (*App, error) {
	a := &App{}
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
	return func(a *App) error {
		a.repoPeople = repo
		return nil
	}
}

func WithAppPolicies(appPolices roles.AppPolices) AppConfiguration {
	return func(a *App) error {
		a.appPolices = appPolices
		return nil
	}
}

func WithStateStore(store state.StateStore) AppConfiguration {
	return func(a *App) error {
		a.stateStore = store
		return nil
	}
}

func (a *App) Save(ctx context.Context, p *peoplepb.Person) (err error) {
	ctx, span := logger.Span(ctx, "Save", p)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(authInfo.Realm, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
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

	err = a.updateStore(ctx, p)
	if err != nil {
		return err
	}

	logger.DebugCtx(ctx).Any("data", p).Msg("Save")
	return nil
}

func (a *App) GetById(ctx context.Context, id string) (person *peoplepb.Person, err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "Save", nil, logger.TraceKVString("person.id", id))
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(authInfo.Realm, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return nil, err
	}

	logger.InfoCtx(ctx).Str("id", id).Msg("GetById")
	person, err = a.repoPeople.GetById(ctx, id)
	if err != nil {
		return
	}

	if person == nil {
		return nil, app.ItemNotFoundError("person", id)
	}

	return person, nil
}

func (a *App) List(ctx context.Context, filters Filters) (result []*peoplepb.Person, err error) {
	ctx, span := logger.Span(ctx, "List", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(authInfo.Realm, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
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

func (a *App) Update(ctx context.Context, p *peoplepb.Person) (err error) {
	ctx, span := logger.Span(ctx, "Update", p)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, peoplepb.PeopleSvc_Save_FullMethodName, authInfo)
	if err != nil {
		return
	}
	logger.DebugCtx(ctx).Any("person", p).Msg("Update")
	err = a.repoPeople.Update(ctx, p)
	if err != nil {
		return
	}

	err = a.updateStore(ctx, p)

	return err
}

func (*App) PinParse(ctx context.Context, number string) (result ddd.PinValidation, err error) {
	_, span := logger.Span(ctx, "PinParse", nil)
	defer func() {
		span.End(err)
	}()

	info, err := pin.Parse(number)
	if err != nil {
		return
	}
	return info, nil
}

func parseEGN(person *peoplepb.Person) {
	res, err := pin.Parse(person.IdNumbers["EGN"])
	if err != nil {
		logger.Error(err).Any("EGN", person)
		return
	}
	if person.Dob == nil {
		person.Dob = &peoplepb.Dob{}
	}
	if res.DOB.Year > 0 {
		person.Dob.Year = uint32(res.DOB.Year)
	}
	if res.DOB.Month > 0 {
		person.Dob.Month = uint32(res.DOB.Month)
	}
	if res.DOB.Day > 0 {
		person.Dob.Day = uint32(res.DOB.Day)
	}
	if len(res.Gender) > 0 {
		person.Gender = res.Gender
	}
	logger.Debug().Any("pin", person)
}

func (a *App) updateStore(ctx context.Context, person *peoplepb.Person) error {
	value, err := proto.Marshal(person)
	if err != nil {
		return fmt.Errorf("updateStore.Marshal: %w", err)
	}

	_, err = a.stateStore.Push(ctx, state.NewEntity{
		Key:   person.Id,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("updateStore.Push: %w", err)
	}

	return nil
}
