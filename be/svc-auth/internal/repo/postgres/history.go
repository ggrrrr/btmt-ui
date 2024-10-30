package postgres

import (
	"context"

	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

func (r *authRepo) SaveHistory(ctx context.Context, info roles.AuthInfo, method string) (err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.SaveHistory", nil, logger.KVString("subject", info.Subject))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	insert into %s ("id", "subject", "method", "device" )
	values($1, $2, $3, $4)
	`, r.historyTable)
	logger.DebugCtx(ctx).
		Str("email", info.Subject).
		Str("sql", sql).
		Msg("SaveHistory")

	_, err = r.db.ExecContext(ctx, sql,
		info.ID,
		info.Subject,
		method,
		device{info.Device},
		// time.Now(),
	)

	return err
}

func (r *authRepo) ListHistory(ctx context.Context, subject string) (authHistory []ddd.AuthHistory, err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.ListHistory", nil, logger.KVString("subject", subject))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	select "id", "subject", "method", "device", "created_at" from  %s 
	where "user" = $1
	`, r.historyTable)
	logger.DebugCtx(ctx).
		Str("subject", subject).
		Str("sql", sql).Msg("ListHistory")

	rows, err := r.db.QueryContext(ctx, sql, subject)
	if err != nil {
		return authHistory, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		err = rows.Err()
		return
	}

	for rows.Next() {
		var dev device
		var row ddd.AuthHistory

		err = rows.Scan(&row.ID, &row.Subject, &row.Method, &dev, &row.CreatedAt)
		if err != nil {
			return authHistory, err
		}
		row.Device = dev.Device
		authHistory = append(authHistory, row)
	}

	return
}
func (r *authRepo) GetHistory(ctx context.Context, id uuid.UUID) (*ddd.AuthHistory, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "repo.GetHistory", nil, logger.KVString("id", id.String()))
	defer func() {
		span.End(err)
	}()

	sql := r.table(`
	select "id", "subject", "method", "device", "created_at" from  %s 
	where "id" = $1
	`, r.historyTable)
	logger.DebugCtx(ctx).
		Str("id", id.String()).
		Str("sql", sql).Msg("ListHistory")

	rows, err := r.db.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	if rows.Next() {
		var authHistory ddd.AuthHistory
		var dev device

		err = rows.Scan(&authHistory.ID, &authHistory.Subject, &authHistory.Method, &dev, &authHistory.CreatedAt)
		if err != nil {
			return nil, err
		}
		authHistory.Device = dev.Device
		return &authHistory, nil
	}
	return nil, nil
}

func (r *authRepo) DeleteHistory(ctx context.Context, id string) (err error) {
	ctx, span := logger.SpanWithAttributes(ctx, "repo.DeleteHistory", nil, logger.KVString("id", id))
	defer func() {
		span.End(err)
	}()
	sql := r.table(`
	delete from %s 
	where id = $1
	`, r.historyTable)
	logger.DebugCtx(ctx).
		Str("id", id).
		Str("sql", sql).
		Msg("DeleteHistory")

	_, err = r.db.ExecContext(ctx, sql, id)

	return err
}
