package rest

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	peoplepb "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

func RegisterGateway(ctx context.Context, mux *runtime.ServeMux, grpcAddr string) error {
	// const apiRoot = "/"

	// gateway := runtime.NewServeMux()

	// RegisterEmailSvcHandlerClient

	// mux.Connect("", gateway.)

	// authpb.RegisterEmailSvcHandlerFromEndpoint(ctx, mux, "", )

	// err := authpb.RegisterEmailSvcHandlerClient(ctx, mux, client)
	err := peoplepb.RegisterPeopleSvcHandlerFromEndpoint(ctx, mux, grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		return err
	}

	// mount the GRPC gateway
	// mux.Mount(apiRoot, gateway)
	// mux.

	return nil
}
