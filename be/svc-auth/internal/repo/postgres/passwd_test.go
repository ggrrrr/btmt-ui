package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/postgres"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func cfg() postgres.Config {
	return postgres.Config{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
		Prefix:   "dev",
	}
}

func TestSaveGetList(t *testing.T) {
	ctx := context.Background()

	db, err := postgres.Connect(cfg())
	require.NoError(t, err)

	testRepo := authRepo{
		otelTracer:  tracer.Tracer(otelScope),
		db:          db,
		passwdTable: "test_passwd",
	}

	_, err = testRepo.db.Exec(testRepo.table(`drop table if exists %s`, testRepo.passwdTable))
	require.NoError(t, err)
	err = create(db, testRepo.table(createPasswdTable, testRepo.passwdTable))
	require.NoError(t, err)

	ts := time.Now()
	testData := ddd.AuthPasswd{
		Subject:     "emai@asd.com",
		Passwd:      "pass1",
		Status:      "stat1",
		RealmRoles:  map[string][]string{"localhost": {"admin"}},
		SystemRoles: []string{"systemRoleAdmin"},
	}
	err = testRepo.SavePasswd(ctx, testData)
	require.NoError(t, err)

	rows, err := testRepo.GetPasswd(ctx, testData.Subject)
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.WithinDuration(t, rows[0].CreatedAt, ts, 1*time.Second)
	testData.CreatedAt = rows[0].CreatedAt

	fmt.Printf("----- %+v", rows[0])

	require.Equal(t, testData, rows[0])

	rows, err = testRepo.ListPasswd(ctx, nil)
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.WithinDuration(t, rows[0].CreatedAt, ts, 1*time.Second)
	testData.CreatedAt = rows[0].CreatedAt
	require.Equal(t, testData, rows[0])

	testData1 := ddd.AuthPasswd{Subject: "email2", Passwd: "pass1", Status: "stat1", SystemRoles: []string{"r1"}}
	err = testRepo.SavePasswd(ctx, testData1)
	require.NoError(t, err)
	rows, err = testRepo.ListPasswd(ctx, nil)
	require.NoError(t, err)
	require.True(t, len(rows) == 2)

}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	db, err := postgres.Connect(cfg())
	require.NoError(t, err)

	testRepo, err := Init(db)
	require.NoError(t, err)

	_, err = testRepo.db.Exec(testRepo.table(`drop table if exists %s`, testRepo.passwdTable))
	require.NoError(t, err)
	err = create(db, testRepo.table(createPasswdTable, testRepo.passwdTable))
	require.NoError(t, err)

	ts := time.Now()
	testData := ddd.AuthPasswd{Subject: "email11", Passwd: "pass1", Status: "stat1", SystemRoles: []string{"r1"}}
	err = testRepo.SavePasswd(ctx, testData)
	require.NoError(t, err)

	rows, err := testRepo.GetPasswd(ctx, "email11")
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.WithinDuration(t, rows[0].CreatedAt, ts, 1*time.Second)
	testData.CreatedAt = rows[0].CreatedAt
	require.Equal(t, testData, rows[0])

	err = testRepo.UpdateStatus(ctx, testData.Subject, "ok")
	require.NoError(t, err)
	rows, err = testRepo.GetPasswd(ctx, testData.Subject)
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.Equal(t, ddd.StatusType("ok"), rows[0].Status)

	err = testRepo.UpdatePassword(ctx, testData.Subject, "asdqweasdqwe")
	require.NoError(t, err)
	rows, err = testRepo.GetPasswd(ctx, testData.Subject)
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.Equal(t, "asdqweasdqwe", rows[0].Passwd)

	updateData := ddd.AuthPasswd{
		Subject:     testData.Subject,
		Status:      ddd.StatusPending,
		SystemRoles: []string{"notadmin", "other"},
		RealmRoles:  map[string][]string{"t1": {"asd"}},
	}
	err = testRepo.Update(ctx, updateData)
	require.NoError(t, err)
	rows, err = testRepo.GetPasswd(ctx, testData.Subject)
	require.NoError(t, err)
	assert.True(t, len(rows) == 1)
	got := rows[0]
	updateData.CreatedAt = got.CreatedAt
	updateData.Passwd = got.Passwd
	assert.Equal(t, updateData, got)

}
