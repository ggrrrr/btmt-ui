package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb"
)

type server struct {
	app app.App
	peoplepb.UnimplementedPeopleSvcServer
}

func RegisterServer(app app.App, registrar grpc.ServiceRegistrar) {
	logger.Info().Msg("grpc.RegisterServer")
	peoplepb.RegisterPeopleSvcServer(registrar, &server{
		app: app,
	})
}

func (s *server) Save(ctx context.Context, req *peoplepb.SaveRequest) (*peoplepb.SaveResponse, error) {
	person := req.ToPerson()
	err := s.app.Save(ctx, person)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("person", &req).Msg("Save")
		return nil, err
	}
	return &peoplepb.SaveResponse{
		Id: person.Id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *peoplepb.GetRequest) (*peoplepb.GetResponse, error) {
	res, err := s.app.GetById(ctx, req.Id)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("req", &req).Msg("Get")
		return nil, err
	}
	return &peoplepb.GetResponse{
		Data: peoplepb.FromPerson(res),
	}, nil
}

func (s *server) List(ctx context.Context, req *peoplepb.ListRequest) (*peoplepb.ListResponse, error) {
	list, err := s.app.List(ctx, req.ToFilter())
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("req", &req).Msg("List")
		return nil, err
	}
	out := peoplepb.ListResponse{}
	if len(list) == 0 {
		logger.ErrorCtx(ctx, err).Any("req", &req).Msg("List")
		return &out, nil
	}
	out.List = []*peoplepb.Person{}
	for _, p := range list {
		out.List = append(out.List, peoplepb.FromPerson(&p))
	}
	return &out, nil
}

func (s *server) Update(ctx context.Context, req *peoplepb.UpdateRequest) (*peoplepb.UpdateResponse, error) {
	p := req.ToPerson()
	err := s.app.Update(ctx, p)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("req", &req).Msg("Update")
		return nil, err
	}
	return &peoplepb.UpdateResponse{}, nil
}
