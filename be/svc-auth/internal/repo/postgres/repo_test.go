package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/postgres"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cfg() postgres.Config {
	return postgres.Config{
		Host:     "localhost",
		Port:     5432,
		Username: "initexample",
		Password: "initexample",
		Database: "test",
		Preffix:  "dev",
	}
}

func TestSave(t *testing.T) {
	ctx := context.Background()
	conn, err := Connect(cfg())
	require.NoError(t, err)

	_, err = conn.db.Exec(conn.table(`delete from %s`))
	require.NoError(t, err)

	ts := time.Now()
	testData := ddd.AuthPasswd{Email: "email11", Passwd: "pass1", Status: "stat1", SystemRoles: []string{"r1"}}
	err = conn.Save(ctx, testData)
	require.NoError(t, err)

	rows, err := conn.Get(ctx, "email11")
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.WithinDuration(t, rows[0].CreatedAt, ts, 1*time.Second)
	testData.CreatedAt = rows[0].CreatedAt
	require.Equal(t, testData, rows[0])

	rows, err = conn.List(ctx)
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.WithinDuration(t, rows[0].CreatedAt, ts, 1*time.Second)
	testData.CreatedAt = rows[0].CreatedAt
	require.Equal(t, testData, rows[0])

	testData1 := ddd.AuthPasswd{Email: "email2", Passwd: "pass1", Status: "stat1", SystemRoles: []string{"r1"}}
	err = conn.Save(ctx, testData1)
	require.NoError(t, err)
	rows, err = conn.List(ctx)
	require.NoError(t, err)
	require.True(t, len(rows) == 2)

}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	conn, err := Connect(cfg())
	require.NoError(t, err)

	_, err = conn.db.Exec(conn.table(`delete from %s`))
	require.NoError(t, err)

	ts := time.Now()
	testData := ddd.AuthPasswd{Email: "email11", Passwd: "pass1", Status: "stat1", SystemRoles: []string{"r1"}}
	err = conn.Save(ctx, testData)
	require.NoError(t, err)

	rows, err := conn.Get(ctx, "email11")
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.WithinDuration(t, rows[0].CreatedAt, ts, 1*time.Second)
	testData.CreatedAt = rows[0].CreatedAt
	require.Equal(t, testData, rows[0])

	err = conn.UpdateStatus(ctx, testData.Email, "ok")
	require.NoError(t, err)
	rows, err = conn.Get(ctx, testData.Email)
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.Equal(t, ddd.StatusType("ok"), rows[0].Status)

	err = conn.UpdatePassword(ctx, testData.Email, "asdqweasdqwe")
	require.NoError(t, err)
	rows, err = conn.Get(ctx, testData.Email)
	require.NoError(t, err)
	require.True(t, len(rows) == 1)
	require.Equal(t, "asdqweasdqwe", rows[0].Passwd)

	updateData := ddd.AuthPasswd{
		Email:       testData.Email,
		Status:      ddd.StatusPending,
		SystemRoles: []string{"notadmin", "other"},
	}
	err = conn.Update(ctx, updateData)
	require.NoError(t, err)
	rows, err = conn.Get(ctx, testData.Email)
	require.NoError(t, err)
	assert.True(t, len(rows) == 1)
	assert.Equal(t, updateData.Status, rows[0].Status)
	assert.Equal(t, updateData.SystemRoles, rows[0].SystemRoles)

}
