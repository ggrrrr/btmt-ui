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
	cfg := mgo.MgoTestCfg("test-temp")
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
		Body:        `some shit asdasd <asd></asd>`,
		CreatedAt:   time.Now(),
	}
	err = testRepo.Save(ctx, firstTmpl)
	require.NoError(t, err)
	assert.True(t, firstTmpl.Id != "")

	actualTmpl, err := testRepo.GetById(ctx, firstTmpl.Id)
	require.NoError(t, err)
	fmt.Printf("  %+v \n", actualTmpl)
	require.NotNil(t, actualTmpl)
	assert.WithinDuration(t, firstTmpl.CreatedAt, actualTmpl.CreatedAt, 100*time.Millisecond)
	firstTmpl.CreatedAt = actualTmpl.CreatedAt
	assert.Equal(t, firstTmpl, actualTmpl)

	listResult, err := testRepo.List(ctx, nil)
	require.NoError(t, err)

	assert.WithinDuration(t, firstTmpl.CreatedAt, listResult[0].CreatedAt, 100*time.Millisecond)
	firstTmpl.CreatedAt = actualTmpl.CreatedAt
	assert.Equal(t, firstTmpl, actualTmpl)
}
