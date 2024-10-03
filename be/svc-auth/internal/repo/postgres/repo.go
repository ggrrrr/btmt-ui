package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

const createTable string = `CREATE TABLE IF NOT EXISTS %s (
		email TEXT  not null,
		passwd TEXT  not null,
		"status" TEXT  not null,
		system_roles TEXT[] not null,
		tenant_roles JSONB  not null,
		created_at TIMESTAMP DEFAULT NOW(),
		UNIQUE(email)
	)`

type (
	tRoles map[string][]string
	repo   struct {
		prefix string
		db     *sql.DB
	}
)

var _ (ddd.AuthPasswdRepo) = (*repo)(nil)

func (r *repo) table(sql string) string {
	return fmt.Sprintf(sql, r.prefix)
}

func Init(db *sql.DB) (*repo, error) {
	r := &repo{
		prefix: "auth",
		db:     db,
	}
	err := create(r.db, r.table(createTable))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func create(db *sql.DB, sql string) error {
	logger.Info().Str("sql", sql).Msg("create table")
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) Get(ctx context.Context, email string) (result []ddd.AuthPasswd, err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.Get", nil, logger.AttributeString("email", email))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	select "email", "passwd", "status", "tenant_roles", "system_roles", created_at from %s
	where email = $1
	`)
	logger.DebugCtx(ctx).
		Str("email", email).
		Str("sql", sql).Msg("Get")
	rows, err := r.db.QueryContext(ctx, sql, email)
	if err != err {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m tRoles
		var row ddd.AuthPasswd
		err = rows.Scan(
			&row.Email,
			&row.Passwd,
			&row.Status,
			&m,
			pq.Array(&row.SystemRoles),
			&row.CreatedAt,
		)
		if err != nil {
			return
		}
		row.TenantRoles = m
		result = append(result, row)
	}
	return
}

func (r *repo) List(ctx context.Context, filter app.FilterFactory) (result []ddd.AuthPasswd, err error) {
	ctx, span := logger.Span(ctx, "repo.List", nil)
	defer func() {
		span.End(err)
	}()

	sql := r.table(`select "email", "passwd", "status", "tenant_roles", "system_roles", created_at from %s`)
	logger.DebugCtx(ctx).
		Str("sql", sql).Msg("List")
	// rows, err := r.db.Query(sql)
	rows, err := r.db.QueryContext(ctx, sql)
	if err != err {
		return
	}
	if rows == nil {
		err = app.SystemError("db.Rows is nil", nil)
		return
	}

	defer rows.Close()

	if rows.Err() != nil {
		return
	}
	out := []ddd.AuthPasswd{}
	for rows.Next() {
		var m tRoles
		var row ddd.AuthPasswd
		err = rows.Scan(
			&row.Email,
			&row.Passwd,
			&row.Status,
			&m,
			pq.Array(&row.SystemRoles),
			&row.CreatedAt,
		)
		if err != nil {
			return
		}
		row.TenantRoles = m

		fmt.Printf("\t\t %v \n", row)

		out = append(out, row)
	}
	return out, nil
}

func (r *repo) Update(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.Update", nil, logger.AttributeString("email", auth.Email))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	update %s set  "status" = $2, "system_roles" = $3, "tenant_roles" = $4
	where email = $1
	`)
	logger.DebugCtx(ctx).
		Str("email", auth.Email).
		Str("sql", sql).Msg("Update")
	_, err = r.db.ExecContext(ctx, sql,
		auth.Email,
		auth.Status,
		pq.Array(auth.SystemRoles),
		tRoles(auth.TenantRoles),
	)
	if err != nil {
		return
	}
	return nil
}

func (r *repo) Save(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.Save", nil, logger.AttributeString("email", auth.Email))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	insert into %s ("email", "passwd", "status", "tenant_roles", "system_roles")
	values($1, $2, $3, $4, $5)
	`)
	logger.DebugCtx(ctx).
		Str("email", auth.Email).
		Str("sql", sql).Msg("Save")
	_, err = r.db.ExecContext(ctx, sql,
		auth.Email,
		auth.Passwd,
		auth.Status,
		tRoles(auth.TenantRoles),
		pq.Array(auth.SystemRoles),
	)
	if err != nil {
		return
	}
	return nil
}

func (r *repo) UpdatePassword(ctx context.Context, email string, password string) (err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.UpdatePassword", nil, logger.AttributeString("email", email))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	update  %s set "passwd" = $1
	where email = $2
	`)
	logger.DebugCtx(ctx).
		Str("email", email).
		Str("sql", sql).Msg("UpdatePassword")
	_, err = r.db.ExecContext(ctx, sql,
		password,
		email,
	)
	if err != nil {
		return
	}
	return nil
}

func (r *repo) UpdateStatus(ctx context.Context, email string, status ddd.StatusType) (err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.UpdateStatus", nil, logger.AttributeString("email", email))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	update  %s set "status" = $1
	where email = $2
	`)
	logger.DebugCtx(ctx).
		Str("email", email).
		Str("sql", sql).Msg("UpdateStatus")
	_, err = r.db.ExecContext(ctx, sql,
		status,
		email,
	)
	if err != nil {
		return
	}
	return
}
