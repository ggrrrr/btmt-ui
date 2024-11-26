package grpc

import (
	"google.golang.org/grpc"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	authpb "github.com/ggrrrr/btmt-ui/be/svc-auth/authpb/v1"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
)

type server struct {
	app app.App
	authpb.UnimplementedAuthSvcServer
}

func RegisterServer(app app.App, registrar grpc.ServiceRegistrar) {
	logger.Info().Msg("grpc.RegisterServer")
	authpb.RegisterAuthSvcServer(registrar, &server{
		app: app,
	})
}
