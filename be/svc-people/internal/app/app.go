package app

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/pin"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	peoplepb "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.be.svc-people"

type (
	OptionFn func(a *Application) error

	Filters map[string][]string

	Application struct {
		tracer     tracer.OTelTracer
		repoPeople ddd.PeopleRepo
		appPolices roles.AppPolices
		stateStore state.StateStore
	}

	App interface {
		Save(ctx context.Context, p *peoplepb.Person) error
		List(ctx context.Context, filters Filters) ([]*peoplepb.Person, error)
		GetById(ctx context.Context, id string) (*peoplepb.Person, error)
		Update(ctx context.Context, p *peoplepb.Person) error
		IDParse(ctx context.Context, pin *peoplepb.IDParseRequest) (*peoplepb.IDParseResponse, error)
	}
)

var _ (App) = (*Application)(nil)

var (
	FilterTexts  string = "texts"
	FilterPhones string = "phones"
	FilterPINs   string = "pins"
	FilterLabels string = "labels"
	FilterAttrs  string = "attrs"
)

// var _ (App) = (*App)(nil)

func New(cfgs ...OptionFn) (*Application, error) {
	a := &Application{
		tracer: tracer.Tracer(otelScope),
	}
	for _, c := range cfgs {
		err := c(a)
		if err != nil {
			return nil, err
		}
	}
	if a.appPolices == nil {
		log.Log().Warn(nil, "use mock AppPolices")
		a.appPolices = roles.NewAppPolices()
	}
	if a.appPolices == nil ||
		a.repoPeople == nil ||
		a.stateStore == nil {
		return nil, fmt.Errorf("interface is nil")
	}
	return a, nil
}

func WithPeopleRepo(repo ddd.PeopleRepo) OptionFn {
	return func(a *Application) error {
		a.repoPeople = repo
		return nil
	}
}

func WithAppPolicies(appPolices roles.AppPolices) OptionFn {
	return func(a *Application) error {
		a.appPolices = appPolices
		return nil
	}
}

func WithStateStore(store state.StateStore) OptionFn {
	return func(a *Application) error {
		a.stateStore = store
		return nil
	}
}

func (a *Application) Save(ctx context.Context, p *peoplepb.Person) (err error) {
	ctx, span := a.tracer.SpanWithData(ctx, "Save", p)
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

	log.Log().DebugCtx(ctx, "save")
	return nil
}

func (a *Application) GetById(ctx context.Context, id string) (person *peoplepb.Person, err error) {
	ctx, span := a.tracer.SpanWithAttributes(ctx, "GetById", slog.String("id", id))
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(authInfo.Realm, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return nil, err
	}

	log.Log().InfoCtx(ctx, "getByID", slog.String("id", id))
	person, err = a.repoPeople.GetById(ctx, id)
	if err != nil {
		return
	}

	if person == nil {
		return nil, app.ItemNotFoundError("person", id)
	}

	return person, nil
}

func (a *Application) List(ctx context.Context, filters Filters) (result []*peoplepb.Person, err error) {
	ctx, span := a.tracer.Span(ctx, "List")
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	if err := a.appPolices.CanDo(authInfo.Realm, peoplepb.PeopleSvc_Save_FullMethodName, authInfo); err != nil {
		return nil, err
	}

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

		return nil, err
	}
	return out, nil
}

func (a *Application) Update(ctx context.Context, p *peoplepb.Person) (err error) {
	ctx, span := a.tracer.SpanWithData(ctx, "Update", p)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, peoplepb.PeopleSvc_Save_FullMethodName, authInfo)
	if err != nil {
		return
	}
	err = a.repoPeople.Update(ctx, p)
	if err != nil {
		return
	}

	err = a.updateStore(ctx, p)

	return err
}

func (*Application) IDParse(ctx context.Context, request *peoplepb.IDParseRequest) (result *peoplepb.IDParseResponse, err error) {

	info, err := pin.Parse(request.Number)
	if err != nil {
		return
	}
	return &peoplepb.IDParseResponse{
		Payload: &peoplepb.IDPayload{
			Dob: &peoplepb.Dob{
				Year:  uint32(info.DOB.Year),
				Month: uint32(info.DOB.Month),
				Day:   uint32(info.DOB.Day),
			},
			Gender: info.Gender,
		},
	}, nil
}

func parseEGN(person *peoplepb.Person) {
	res, err := pin.Parse(person.IdNumbers["EGN"])
	if err != nil {
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
}

func (a *Application) updateStore(ctx context.Context, person *peoplepb.Person) error {
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
