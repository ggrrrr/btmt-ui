package repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

func TestSave(t *testing.T) {
	collection := "test-templ"
	var err error
	ctx := context.Background()
	testDb, err := mgo.ConnectForTest(collection)
	require.NoError(t, err)
	// defer testRepo.Close()
	defer testDb.Close(ctx)

	testRepo := New(collection, testDb)

	err = testDb.DB().Collection(testRepo.collection).Drop(ctx)
	require.NoError(t, err)

	// err = testDb.DB().CreateCollection(ctx, testRepo.collection)
	// require.NoError(t, err)

	firstTmpl := &tmplpb.Template{
		Labels:      []string{"label1"},
		Name:        "test template",
		ContentType: "ctx/html",
		Body: `<p> {{ .UserInfo.Device.DeviceInfo }}</p>
<p> {{ .UserInfo.Subject }}</p>
{{ renderImg "IMG4944.JPG:1" }}`,
	}
	err = testRepo.Save(ctx, firstTmpl)
	require.NoError(t, err)
	assert.True(t, firstTmpl.Id != "")

	actualTmpl, err := testRepo.GetById(ctx, firstTmpl.Id)
	require.NoError(t, err)
	require.NotNil(t, actualTmpl)

	tmplpb.MatchTemplate(t, time.Now(), firstTmpl, actualTmpl)

	listResult, err := testRepo.List(ctx, nil)
	require.NoError(t, err)

	tmplpb.MatchTemplate(t, time.Now(), firstTmpl, listResult[0])

	updateTmpl := &tmplpb.Template{
		Id:          firstTmpl.Id,
		Labels:      []string{"label1", "label2"},
		Name:        "updated template",
		ContentType: "ctx/html",
		Body: `<p>From update</p>
<p> {{ .UserInfo.Device.DeviceInfo }}</p>
<p> {{ .UserInfo.Subject }}</p>
{{ renderImg "IMG4944.JPG:1" }}`,
		CreatedAt: firstTmpl.CreatedAt,
		UpdatedAt: timestamppb.Now(),
	}

	err = testRepo.Update(ctx, updateTmpl)
	require.NoError(t, err)
	actualTmpl, err = testRepo.GetById(ctx, firstTmpl.Id)
	require.NoError(t, err)

	require.NotNil(t, actualTmpl)
	tmplpb.MatchTemplate(t, time.Now(), updateTmpl, actualTmpl)

	actualTmpl, err = testRepo.GetById(ctx, updateTmpl.Id)
	require.NoError(t, err)

	tmplpb.MatchTemplate(t, time.Now(), updateTmpl, actualTmpl)

}
