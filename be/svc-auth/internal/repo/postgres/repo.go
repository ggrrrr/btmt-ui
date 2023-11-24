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
	repo struct {
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

func Connect(cfg postgres.Config) (*repo, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if cfg.Preffix == "" {
		cfg.Preffix = "dev"
	}
	logger.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("User", cfg.Username).
		Msg("Connected")
	return &repo{
		db:     db,
		prefix: strings.Trim(cfg.Preffix, " "),
	}, nil
}

func (r *repo) Get(ctx context.Context, email string) ([]ddd.AuthPasswd, error) {
	sql := r.table(`
	select "email", "passwd", "status", "system_roles", created_at from %s_auth
	where email = $1
	`)
	logger.DebugCtx(ctx).Str("sql", sql).Send()
	rows, err := r.db.QueryContext(ctx, sql, email)
	if err != err {
		return []ddd.AuthPasswd{}, err
	}
	out := []ddd.AuthPasswd{}
	for rows.Next() {
		var row ddd.AuthPasswd
		err = rows.Scan(
			&row.Email,
			&row.Passwd,
			&row.Status,
			pq.Array(&row.SystemRoles),
			&row.CreatedAt,
		)
		if err != nil {
			return []ddd.AuthPasswd{}, err
		}
		out = append(out, row)
	}
	return out, nil
}

func (r *repo) List(ctx context.Context) ([]ddd.AuthPasswd, error) {
	sql := r.table(`
	select "email", "passwd", "status", "system_roles", created_at from %s_auth
	`)
	logger.DebugCtx(ctx).Str("sql", sql).Send()
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
		var row ddd.AuthPasswd
		err = rows.Scan(
			&row.Email,
			&row.Passwd,
			&row.Status,
			pq.Array(&row.SystemRoles),
			&row.CreatedAt,
		)
		if err != nil {
			return []ddd.AuthPasswd{}, err
		}
		out = append(out, row)
	}
	return out, nil
}

func (r *repo) Save(ctx context.Context, auth ddd.AuthPasswd) error {
	sql := r.table(`
	insert into %s_auth("email", "passwd", "status", "system_roles")
	values($1, $2, $3, $4) 
	`)
	logger.DebugCtx(ctx).Str("sql", sql).Send()
	_, err := r.db.ExecContext(ctx, sql,
		auth.Email,
		auth.Passwd,
		auth.Status,
		pq.Array(auth.SystemRoles),
	)
	return err
}

func (r *repo) UpdatePassword(ctx context.Context, email string, password string) error {
	sql := r.table(`
	update  %s_auth set "passwd" = $1
	where email = $2
	`)
	logger.DebugCtx(ctx).Str("sql", sql).Send()
	_, err := r.db.ExecContext(ctx, sql,
		password,
		email,
	)
	return err
}

func (r *repo) UpdateStatus(ctx context.Context, email string, status ddd.StatusType) error {
	sql := r.table(`
	update  %s_auth set "status" = $1
	where email = $2
	`)
	logger.DebugCtx(ctx).Str("sql", sql).Send()
	_, err := r.db.ExecContext(ctx, sql,
		status,
		email,
	)
	return err
}
