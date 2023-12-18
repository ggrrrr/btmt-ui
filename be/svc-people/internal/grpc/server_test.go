package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	"github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb"
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

func TestTelephoneServer_GetContact(t *testing.T) {
	ctx := context.Background()

	cfg := mgo.MgoTestCfg()

	testDb, err := mgo.New(ctx, cfg)
	require.NoError(t, err)

	testRepo := repo.New(cfg.Collection, testDb)
	require.NoError(t, err)
	// defer testRepo.Close()

	testApp, err := app.New(
		app.WithPeopleRepo(testRepo),
		app.WithAppPolicies(roles.NewAppPolices()),
	)
	require.NoError(t, err)

	client, closer := testServer(ctx, testApp)
	defer closer()

	tests := []testCase{
		{
			test: "list ok",
			testFunc: func(tt *testing.T) {
				res, err := client.List(ctx, &peoplepb.ListRequest{Filters: map[string]*peoplepb.ListText{}})
				assert.NoError(tt, err)
				fmt.Printf("%v \n", res)
				//
			},
		},
		{
			test: "save get ok",
			testFunc: func(tt *testing.T) {
				res, err := client.Save(ctx, &peoplepb.SaveRequest{
					Data: &peoplepb.Person{
						IdNumbers: map[string]string{"pin": "pin1"},
						Name:      "save ok test",
						Phones:    map[string]string{"mobile": "123123123"},
					},
				})
				require.NoError(tt, err)
				fmt.Printf("%v \n", res)
				res1, err := client.Get(ctx, &peoplepb.GetRequest{Id: res.Id})
				require.NoError(tt, err)
				fmt.Printf("%v \n", res1)
				//
				res1.Data.FullName = "some full name"
				res2, err := client.Update(ctx, &peoplepb.UpdateRequest{
					Data: res1.Data,
				})
				require.NoError(tt, err)
				fmt.Printf("%v \n", res2)
				res3, err := client.Get(ctx, &peoplepb.GetRequest{Id: res.Id})
				require.NoError(tt, err)
				fmt.Printf("%v \n", res3)
				assert.Equal(tt, res3.Data.FullName, "some full name")

			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.test, tc.testFunc)
	}

}

func testServer(ctx context.Context, app app.App) (peoplepb.PeopleSvcClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			ctx = roles.CtxWithAuthInfo(ctx, roles.CreateAdminUser("test", roles.Device{}))
			return handler(ctx, req)
		}),
	)
	peoplepb.RegisterPeopleSvcServer(baseServer, &server{app: app})
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
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

	client := peoplepb.NewPeopleSvcClient(conn)

	return client, closer
}
