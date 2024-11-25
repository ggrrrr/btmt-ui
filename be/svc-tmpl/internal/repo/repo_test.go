package repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func TestSave(t *testing.T) {
	var err error
	ctx := context.Background()
	cfg := mgo.MgoTestCfg("tmpl")
	testDb, err := mgo.New(ctx, cfg)
	require.NoError(t, err)
	// defer testRepo.Close()
	defer testDb.Close(ctx)

	testRepo := New(cfg.Collection, testDb)

	err = testDb.DB().Collection(testRepo.collection).Drop(ctx)
	require.NoError(t, err)

	// err = testDb.DB().CreateCollection(ctx, testRepo.collection)
	// require.NoError(t, err)

	firstTmpl := &ddd.Template{
		Labels:      []string{"label1"},
		Name:        "test template",
		ContentType: "ctx/html",
		Body: `<p> {{ .UserInfo.Device.DeviceInfo }}</p>
<p> {{ .UserInfo.Subject }}</p>
{{ renderImg "IMG4944.JPG:1" }}`,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = testRepo.Save(ctx, firstTmpl)
	require.NoError(t, err)
	assert.True(t, firstTmpl.Id != "")

	actualTmpl, err := testRepo.GetById(ctx, firstTmpl.Id)
	require.NoError(t, err)
	require.NotNil(t, actualTmpl)

	ddd.MatchTemplate(t, *firstTmpl, *actualTmpl)

	listResult, err := testRepo.List(ctx, nil)
	require.NoError(t, err)

	ddd.MatchTemplate(t, *firstTmpl, listResult[0])

	updateTmpl := &ddd.Template{
		Id:          firstTmpl.Id,
		Labels:      []string{"label1", "label2"},
		Name:        "updated template",
		ContentType: "ctx/html",
		Body: `<p>From update</p>
<p> {{ .UserInfo.Device.DeviceInfo }}</p>
<p> {{ .UserInfo.Subject }}</p>
{{ renderImg "IMG4944.JPG:1" }}`,
		CreatedAt: firstTmpl.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = testRepo.Update(ctx, updateTmpl)
	require.NoError(t, err)
	actualTmpl, err = testRepo.GetById(ctx, firstTmpl.Id)
	require.NoError(t, err)

	require.NotNil(t, actualTmpl)
	assert.WithinDuration(t, updateTmpl.UpdatedAt, actualTmpl.UpdatedAt, 100*time.Millisecond)
	ddd.MatchTemplate(t, *updateTmpl, *actualTmpl)

	updateTmpl.BlobId = "blob_id_11"
	err = testRepo.UpdateBlobId(ctx, updateTmpl)
	require.NoError(t, err)

	actualTmpl, err = testRepo.GetById(ctx, updateTmpl.Id)
	require.NoError(t, err)

	ddd.MatchTemplate(t, *updateTmpl, *actualTmpl)

}
