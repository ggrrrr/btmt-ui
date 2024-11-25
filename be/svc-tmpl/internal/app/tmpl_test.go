package app

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/awsclient"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/blob/awss3"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/repo"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb"
)

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

	blobClient, err := awss3.NewClient("test-bucket-1", awsclient.AwsConfig{
		Region:   "us-east-1",
		Endpoint: "http://localhost:4566",
	})
	require.NoError(t, err)

	testApp, err := New(WithBlobStore(blobClient), WithTmplRepo(testRepo))
	require.NoError(t, err)

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "save new",
			testFunc: func(t *testing.T) {
				var err error
				newTmpl := &tmplpb.Template{
					Id:          "",
					ContentType: "text/html",
					Name:        "test tmpl",
					Labels:      []string{"new:label"},
					Images:      []string{"image-1"},
					Files:       map[string]string{},
					Body:        "new template body",
					// CreatedAt:   ,
					// UpdatedAt:   time.Time{},
				}

				tmplErr, err := testApp.SaveTmpl(ctx, newTmpl)
				require.NoError(t, err)
				require.Equal(t, TmplError(nil), tmplErr)
				require.True(t, newTmpl.Id != "")

				savedTmpl, err := testApp.GetTmpl(ctx, newTmpl.Id)
				require.NoError(t, err)

				tmplpb.MatchTemplate(t, newTmpl, savedTmpl)

				blobId1, err := testApp.tmplFolder.SetIdVersionFromString(savedTmpl.Id)
				require.NoError(t, err)

				blobId, err := testApp.tmplFolder.SetIdVersionFromString(blobId1.Id())
				require.NoError(t, err)

				// fmt.Printf("blobId: %#v, \n", blobId)

				blobList, err := blobClient.List(ctx, realm, blobId)
				require.NoError(t, err)
				require.Equal(t, len(blobList), 1)
				require.Equal(t, len(blobList[0].Versions), 0)

				require.Equal(t, blobList[0].Id.Id(), savedTmpl.Id)
				require.Equal(t, blobList[0].MD.ContentType, savedTmpl.ContentType)
				require.Equal(t, blobList[0].MD.Type, blob.BlobTypeTemplate)

				tmplBlob, err := blobClient.Fetch(ctx, realm, blobId)
				require.NoError(t, err)

				tmplBlobBody, err := io.ReadAll(tmplBlob.ReadCloser)
				require.NoError(t, err)
				require.Equal(t, newTmpl.Body, string(tmplBlobBody))

				fmt.Printf("blobList: %v, \n", string(tmplBlobBody))

				savedTmpl1, err := testApp.GetTmpl(ctx, newTmpl.Id)
				require.NoError(t, err)

				savedTmpl1.Labels = []string{"new:update-label"}
				_, err = testApp.SaveTmpl(ctx, savedTmpl1)
				require.NoError(t, err)

				savedTmpl1Actual, err := testApp.GetTmpl(ctx, savedTmpl1.Id)
				require.NoError(t, err)
				tmplpb.MatchTemplate(t, savedTmpl1, savedTmpl1Actual)

				blobList, err = blobClient.List(ctx, realm, blobId)
				require.NoError(t, err)
				require.Equal(t, len(blobList), 1)
				require.Equal(t, len(blobList[0].Versions), 0)

				savedTmpl1.Body = "update body from test"
				savedTmpl1.UpdatedAt = timestamppb.New(time.Now())
				_, err = testApp.SaveTmpl(ctx, savedTmpl1)
				require.NoError(t, err)

				savedTmpl1Actual, err = testApp.GetTmpl(ctx, savedTmpl1.Id)
				require.NoError(t, err)
				tmplpb.MatchTemplate(t, savedTmpl1, savedTmpl1Actual)

				fmt.Printf("size:  %#v \n\n\n", savedTmpl1Actual)

				blobList, err = blobClient.List(ctx, realm, blobId)
				require.NoError(t, err)
				require.Equal(t, len(blobList), 1)
				require.Equal(t, len(blobList[0].Versions), 1)

				fmt.Printf("size: %v   %#v \n", len(blobList), blobList)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
