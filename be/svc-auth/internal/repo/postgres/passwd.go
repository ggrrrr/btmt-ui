package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lib/pq"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (r *authRepo) GetPasswd(ctx context.Context, subject string) (result []ddd.AuthPasswd, err error) {
	ctx, span := r.otelTracer.SpanWithAttributes(
		ctx, "repo.GetPasswd",
		slog.String("subject", subject))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	select "subject", "passwd", "status", "tenant_roles", "system_roles", created_at from %s
	where subject = $1
	`, r.passwdTable)
	rows, err := r.db.QueryContext(ctx, sql, subject)
	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Err() != nil {
		err = fmt.Errorf("rows %w", rows.Err())
		return
	}

	if rows == nil {
		err = fmt.Errorf("rows is nil")
		return
	}

	for rows.Next() {
		var m tRoles
		var row ddd.AuthPasswd
		err = rows.Scan(
			&row.Subject,
			&row.Passwd,
			&row.Status,
			&m,
			pq.Array(&row.SystemRoles),
			&row.CreatedAt,
		)
		if err != nil {
			return
		}
		row.RealmRoles = m
		result = append(result, row)
	}
	return
}

func (r *authRepo) ListPasswd(ctx context.Context, filter app.FilterFactory) (result []ddd.AuthPasswd, err error) {
	ctx, span := r.otelTracer.SpanWithAttributes(
		ctx, "repo.ListPasswd")
	defer func() {
		span.End(err)
	}()

	sql := r.table(`select "subject", "passwd", "status", "tenant_roles", "system_roles", created_at from %s`, r.passwdTable)
	rows, err := r.db.QueryContext(ctx, sql)
	if err != nil {
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
			&row.Subject,
			&row.Passwd,
			&row.Status,
			&m,
			pq.Array(&row.SystemRoles),
			&row.CreatedAt,
		)
		if err != nil {
			return
		}
		row.RealmRoles = m

		out = append(out, row)
	}
	return out, nil
}

func (r *authRepo) Update(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	ctx, span := r.otelTracer.SpanWithAttributes(
		ctx, "repo.Update",
		slog.String("Subject", auth.Subject))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	update %s set  "status" = $2, "system_roles" = $3, "tenant_roles" = $4
	where subject = $1
	`, r.passwdTable)
	_, err = r.db.ExecContext(ctx, sql,
		auth.Subject,
		auth.Status,
		pq.Array(auth.SystemRoles),
		tRoles(auth.RealmRoles),
	)
	if err != nil {
		return
	}
	return nil
}

func (r *authRepo) SavePasswd(ctx context.Context, auth ddd.AuthPasswd) (err error) {
	ctx, span := r.otelTracer.SpanWithAttributes(
		ctx, "repo.Update",
		slog.String("subject", auth.Subject))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	insert into %s ("subject", "passwd", "status", "tenant_roles", "system_roles")
	values($1, $2, $3, $4, $5)
	`, r.passwdTable)
	_, err = r.db.ExecContext(ctx, sql,
		auth.Subject,
		auth.Passwd,
		auth.Status,
		tRoles(auth.RealmRoles),
		pq.Array(auth.SystemRoles),
	)
	if err != nil {
		return
	}
	return nil
}

func (r *authRepo) UpdatePassword(ctx context.Context, subject string, password string) (err error) {
	ctx, span := r.otelTracer.SpanWithAttributes(
		ctx, "repo.UpdatePassword",
		slog.String("subject", subject))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	update  %s set "passwd" = $1
	where subject = $2
	`, r.passwdTable)
	_, err = r.db.ExecContext(ctx, sql,
		password,
		subject,
	)
	if err != nil {
		return
	}
	return nil
}

func (r *authRepo) UpdateStatus(ctx context.Context, subject string, status ddd.StatusType) (err error) {
	ctx, span := r.otelTracer.SpanWithAttributes(
		ctx, "repo.UpdateStatus",
		slog.String("subject", subject),
		slog.String("subject", string(status)),
	)
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	update  %s set "status" = $1
	where subject = $2
	`, r.passwdTable)

	_, err = r.db.ExecContext(ctx, sql,
		status,
		subject,
	)
	if err != nil {
		return
	}
	return
}
