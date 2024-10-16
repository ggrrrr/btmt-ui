package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type (
	testCase struct {
		test     string
		testFunc func(tt *testing.T)
	}
)

// func cfg() awsdb.AwsConfig {
// 	return awsdb.AwsConfig{
// 		Region:   "us-east-1",
// 		Endpoint: "http://localhost:4566",
// 		Database: awsdb.DynamodbConfig{
// 			Database: "",
// 			Prefix:   "test",
// 		},
// 	}
// }

func TestServer(t *testing.T) {
	ctx := context.Background()
	ctxAdmin := roles.CtxWithAuthInfo(ctx, roles.CreateSystemAdminUser(roles.SystemRealm, "admin", roles.Device{}))

	store, err := mem.New()
	require.NoError(t, err)

	testApp, err := app.New(app.WithAuthRepo(store), app.WithTokenSigner(token.NewSignerMock()))
	require.NoError(t, err)

	err = testApp.UserCreate(ctxAdmin, ddd.AuthPasswd{
		Email:  "asd@asd",
		Passwd: "asdasdasd",
		Status: ddd.StatusEnabled,
	})
	require.NoError(t, err)

	client, closer := testServer(ctx, testApp)
	defer closer()

	tests := []testCase{
		{
			test: "login ok",
			testFunc: func(tt *testing.T) {
				_, err = client.LoginPasswd(ctx, &authpb.LoginPasswdRequest{
					Email:    "asd@asd",
					Password: "asdasdasd",
				})
				assert.NoError(tt, err)
			},
		},
		{
			test: "login cfg",
			testFunc: func(tt *testing.T) {
				_, err = client.Oauth2Config(ctx, &authpb.Oauth2ConfigRequest{})
				assert.NoError(tt, err)
			},
		},
		{
			test: "list cfg",
			testFunc: func(tt *testing.T) {
				_, err = client.UserList(ctx, &authpb.UserListRequest{})
				assert.NoError(tt, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, tc.testFunc)
	}

}

func testServer(_ context.Context, app app.App) (authpb.AuthSvcClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			ctx = roles.CtxWithAuthInfo(ctx, roles.CreateSystemAdminUser(roles.SystemRealm, "test", roles.Device{}))
			return handler(ctx, req)
		}),
	)
	authpb.RegisterAuthSvcServer(baseServer, &server{app: app})
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := authpb.NewAuthSvcClient(conn)

	return client, closer
}
