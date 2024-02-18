package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/postgres"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

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

func (r *repo) Close() error {
	return r.db.Close()
}

func (r *repo) create() error {
	sql := r.table(`CREATE TABLE IF NOT EXISTS %s (
		email TEXT  not null,
		passwd TEXT  not null,
		"status" TEXT  not null,
		system_roles TEXT[] not null,
		tenant_roles JSONB  not null,
		created_at TIMESTAMP DEFAULT NOW(),
		UNIQUE(email)
	)`)
	logger.Info().Str("sql", sql).Msg("create table")
	_, err := r.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func Connect(cfg postgres.Config) (*repo, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		logger.Error(err).Msg("Connect")
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logger.Error(err).Msg("Ping")
		return nil, err
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "dev"
	}
	logger.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("User", cfg.Username).
		Msg("Connected")
	repo := &repo{
		db:     db,
		prefix: strings.Trim(cfg.Prefix, " "),
	}

	err = repo.create()
	if err != nil {
		db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *repo) Get(ctx context.Context, email string) ([]ddd.AuthPasswd, error) {
	sql := r.table(`
	select "email", "passwd", "status", "tenant_roles", "system_roles", created_at from %s
	where email = $1
	`)
	logger.DebugCtx(ctx).
		Str("email", email).
		Str("sql", sql).Msg("Get")
	rows, err := r.db.QueryContext(ctx, sql, email)
	if err != err {
		return []ddd.AuthPasswd{}, err
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
			return []ddd.AuthPasswd{}, err
		}
		row.TenantRoles = m
		out = append(out, row)
	}
	return out, nil
}

func (r *repo) List(ctx context.Context) ([]ddd.AuthPasswd, error) {
	sql := r.table(`
	select "email", "passwd", "status", "tenant_roles", "system_roles", created_at from %s
	`)
	logger.DebugCtx(ctx).
		Str("sql", sql).Msg("List")
	rows, err := r.db.QueryContext(ctx, sql)
	if err != err {
		return []ddd.AuthPasswd{}, err
	}
	if rows == nil {
		return []ddd.AuthPasswd{}, fmt.Errorf("rows is nil")
	}
	if rows.Err() != nil {
		return []ddd.AuthPasswd{}, rows.Err()
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
			return []ddd.AuthPasswd{}, err
		}
		row.TenantRoles = m
		out = append(out, row)
	}
	return out, nil
}

func (r *repo) Update(ctx context.Context, auth ddd.AuthPasswd) error {
	sql := r.table(`
	update %s set  "status" = $2, "system_roles" = $3, "tenant_roles" = $4
	where email = $1
	`)
	logger.DebugCtx(ctx).
		Str("email", auth.Email).
		Str("sql", sql).Msg("Update")
	_, err := r.db.ExecContext(ctx, sql,
		auth.Email,
		auth.Status,
		pq.Array(auth.SystemRoles),
		tRoles(auth.TenantRoles),
	)
	return err
}

func (r *repo) Save(ctx context.Context, auth ddd.AuthPasswd) error {
	sql := r.table(`
	insert into %s ("email", "passwd", "status", "tenant_roles", "system_roles")
	values($1, $2, $3, $4, $5)
	`)
	logger.DebugCtx(ctx).
		Str("email", auth.Email).
		Str("sql", sql).Msg("Save")
	_, err := r.db.ExecContext(ctx, sql,
		auth.Email,
		auth.Passwd,
		auth.Status,
		tRoles(auth.TenantRoles),
		pq.Array(auth.SystemRoles),
	)
	return err
}

func (r *repo) UpdatePassword(ctx context.Context, email string, password string) error {
	sql := r.table(`
	update  %s set "passwd" = $1
	where email = $2
	`)
	logger.DebugCtx(ctx).
		Str("email", email).
		Str("sql", sql).Msg("UpdatePassword")
	_, err := r.db.ExecContext(ctx, sql,
		password,
		email,
	)
	return err
}

func (r *repo) UpdateStatus(ctx context.Context, email string, status ddd.StatusType) error {
	sql := r.table(`
	update  %s set "status" = $1
	where email = $2
	`)
	logger.DebugCtx(ctx).
		Str("email", email).
		Str("sql", sql).Msg("UpdateStatus")
	_, err := r.db.ExecContext(ctx, sql,
		status,
		email,
	)
	return err
}
