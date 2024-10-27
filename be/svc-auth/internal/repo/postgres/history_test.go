package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/postgres"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistory(t *testing.T) {
	ctx := context.Background()

	db, err := postgres.Connect(cfg())
	require.NoError(t, err)

	testRepo := authRepo{
		db:           db,
		historyTable: "test_history",
	}

	_, err = testRepo.db.Exec(testRepo.table(`drop table if exists %s`, testRepo.historyTable))
	require.NoError(t, err)
	err = create(db, testRepo.table(createHistoryTable, testRepo.historyTable))
	require.NoError(t, err)

	info := roles.AuthInfo{
		ID:    uuid.New(),
		User:  "someshit",
		Realm: "localhost",
		Device: roles.Device{
			RemoteAddr: "127.0.0.1",
			DeviceInfo: "some device os",
		},
	}

	authHistory := ddd.AuthHistory{
		ID:        info.ID,
		User:      info.User,
		Method:    "/some/login",
		Device:    info.Device,
		CreatedAt: time.Now(),
	}

	err = testRepo.SaveHistory(ctx, info, authHistory.Method)
	require.NoError(t, err)

	list, err := testRepo.ListHistory(ctx, info.User)
	require.NoError(t, err)
	assert.Equal(t, 1, len(list))
	require.WithinDuration(t, authHistory.CreatedAt, list[0].CreatedAt, 100*time.Millisecond)
	authHistory.CreatedAt = list[0].CreatedAt
	assert.Equal(t, authHistory, list[0])

	authRecord, err := testRepo.GetHistory(ctx, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, authRecord)

}
