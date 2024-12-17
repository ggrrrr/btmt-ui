package app

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/blob/awss3"
	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/repo"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

var natsCfg = jetstream.Config{
	URL: "localhost:4222",
}

func TestMock(t *testing.T) {
	blobStore := new(blob.MockBlobStore)
	stateStore := new(state.MockStore)

	realm := "localhost"
	ctx := context.Background()
	ctx = roles.CtxWithAuthInfo(ctx, roles.CreateSystemAdminUser(realm, "admin", app.Device{}))
	var err error

	cfg := mgo.MgoTestCfg("tmpl")
	testDb, err := mgo.New(ctx, cfg)
	require.NoError(t, err)
	// defer testRepo.Close()
	defer testDb.Close(ctx)
	testRepo := repo.New(cfg.Collection, testDb)

	_, err = New(WithBlobStore(blobStore), WithTmplRepo(testRepo), WithStateStore(stateStore))
	require.NoError(t, err)

}

func Test_Save(t *testing.T) {

	realm := "localhost"

	var err error
	ctx := context.Background()

	ctx = roles.CtxWithAuthInfo(ctx, roles.CreateSystemAdminUser(realm, "admin", app.Device{}))

	cfg := mgo.MgoTestCfg("tmpl")
	testDb, err := mgo.New(ctx, cfg)
	require.NoError(t, err)
	// defer testRepo.Close()
	defer testDb.Close(ctx)

	testRepo := repo.New(cfg.Collection, testDb)

	stateStore, err := jetstream.NewStateStore(ctx, natsCfg, state.MustParseEntityType("templates"))
	require.NoError(t, err)

	blobClient, err := awss3.NewClient("test-bucket-1", awsclient.AwsConfig{
		Region:   "us-east-1",
		Endpoint: "http://localhost:4566",
	})
	require.NoError(t, err)

	testApp, err := New(WithBlobStore(blobClient), WithTmplRepo(testRepo), WithStateStore(stateStore))
	require.NoError(t, err)

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "save new",
			testFunc: func(t *testing.T) {
				var err error
				newTmpl := &tmplpb.TemplateUpdate{
					Id:          "",
					ContentType: "text/html",
					Name:        "test tmpl",
					Labels:      []string{"new:label"},
					// Images:      []string{"image-1"},
					Body: "new template body",
					// CreatedAt:   ,
					// UpdatedAt:   time.Time{},
				}

				tmplId, err := testApp.SaveTmpl(ctx, newTmpl)
				require.NoError(t, err)
				require.True(t, tmplId != "")

				newTmpl.Id = tmplId
				savedTmpl, err := testApp.GetTmpl(ctx, newTmpl.Id)
				require.NoError(t, err)

				tmplpb.MatchTemplateUpdate(t, time.Now(), []string{}, newTmpl, savedTmpl)

				savedTmpl1, err := testApp.GetTmpl(ctx, newTmpl.Id)
				require.NoError(t, err)

				newTmpl.Id = tmplId
				newTmpl.Labels = []string{"new:update-label"}
				tmplId, err = testApp.SaveTmpl(ctx, newTmpl)
				require.NoError(t, err)
				assert.Equal(t, newTmpl.Id, tmplId)

				savedTmpl1Actual, err := testApp.GetTmpl(ctx, savedTmpl1.Id)
				require.NoError(t, err)
				tmplpb.MatchTemplateUpdate(t, time.Now(), []string{}, newTmpl, savedTmpl1Actual)

				// blobList, err = blobClient.List(ctx, realm, blobId)
				// require.NoError(t, err)
				// require.Equal(t, len(blobList), 1)
				// require.Equal(t, len(blobList[0].Versions), 0)

				newTmpl.Body = "update body from test"
				// savedTmpl1.UpdatedAt = timestamppb.New(time.Now())
				_, err = testApp.SaveTmpl(ctx, newTmpl)
				require.NoError(t, err)

				savedTmpl1Actual, err = testApp.GetTmpl(ctx, savedTmpl1.Id)
				require.NoError(t, err)
				tmplpb.MatchTemplateUpdate(t, time.Now(), []string{}, newTmpl, savedTmpl1Actual)

				fmt.Printf("size:  %#v \n\n\n", savedTmpl1Actual)

				// blobList, err = blobClient.List(ctx, realm, blobId)
				// require.NoError(t, err)
				// require.Equal(t, len(blobList), 1)
				// require.Equal(t, len(blobList[0].Versions), 1)

				// fmt.Printf("size: %v   %#v \n", len(blobList), blobList)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
