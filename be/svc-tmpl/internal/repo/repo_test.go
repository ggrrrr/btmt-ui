package repo

import (
	"context"
	"fmt"
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
	fmt.Printf("  %+v \n", actualTmpl)
	require.NotNil(t, actualTmpl)
	assert.WithinDuration(t, firstTmpl.CreatedAt, actualTmpl.CreatedAt, 100*time.Millisecond)
	assert.WithinDuration(t, firstTmpl.UpdatedAt, actualTmpl.UpdatedAt, 100*time.Millisecond)
	firstTmpl.CreatedAt = actualTmpl.CreatedAt
	firstTmpl.UpdatedAt = actualTmpl.UpdatedAt
	assert.Equal(t, firstTmpl, actualTmpl)

	listResult, err := testRepo.List(ctx, nil)
	require.NoError(t, err)

	assert.WithinDuration(t, firstTmpl.CreatedAt, listResult[0].CreatedAt, 100*time.Millisecond)
	firstTmpl.CreatedAt = actualTmpl.CreatedAt
	assert.Equal(t, firstTmpl, actualTmpl)

	updateTmpl := &ddd.Template{
		Id:          firstTmpl.Id,
		Labels:      []string{"label1", "label2"},
		Name:        "updated template",
		ContentType: "ctx/html",
		Body: `<p>From update</p>
<p> {{ .UserInfo.Device.DeviceInfo }}</p>
<p> {{ .UserInfo.Subject }}</p>
{{ renderImg "IMG4944.JPG:1" }}`,
		UpdatedAt: time.Now(),
	}

	err = testRepo.Update(ctx, updateTmpl)
	require.NoError(t, err)
	actualTmpl, err = testRepo.GetById(ctx, firstTmpl.Id)
	require.NoError(t, err)

	require.NotNil(t, actualTmpl)
	assert.WithinDuration(t, updateTmpl.UpdatedAt, actualTmpl.UpdatedAt, 100*time.Millisecond)
	updateTmpl.CreatedAt = actualTmpl.CreatedAt
	updateTmpl.UpdatedAt = actualTmpl.UpdatedAt
	assert.Equal(t, updateTmpl, actualTmpl)

	updateTmpl.BlobId = "blob_id_11"
	err = testRepo.UpdateBlobId(ctx, updateTmpl)
	require.NoError(t, err)

	actualTmpl, err = testRepo.GetById(ctx, updateTmpl.Id)
	require.NoError(t, err)
	actualTmpl.CreatedAt = updateTmpl.CreatedAt
	actualTmpl.UpdatedAt = updateTmpl.UpdatedAt
	assert.Equal(t, updateTmpl, actualTmpl)

}
