package mem

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

type (
	repo struct {
		mx *sync.Mutex
		db map[string]*ddd.AuthPasswd
		h  []ddd.AuthHistory
	}
)

// DeleteHistory implements ddd.AuthHistoryRepo.
func (r *repo) DeleteHistory(ctx context.Context, id string) (err error) {
	return nil
}

func (r *repo) ListHistory(ctx context.Context, user string) (authHistory []ddd.AuthHistory, err error) {
	return r.h, nil
}

func (r *repo) GetHistory(ctx context.Context, id uuid.UUID) (authHistory *ddd.AuthHistory, err error) {
	if len(r.h) == 0 {
		return authHistory, fmt.Errorf("not found")
	}
	return &r.h[0], nil
}

func (r *repo) SaveHistory(ctx context.Context, info roles.AuthInfo, method string) (err error) {
	r.h = append(r.h, ddd.AuthHistory{ID: info.ID, Subject: info.Subject, Method: method, Device: info.Device})
	return nil
}

var _ (ddd.AuthPasswdRepo) = (*repo)(nil)
var _ (ddd.AuthHistoryRepo) = (*repo)(nil)

func New() (*repo, error) {
	logger.Warn().Msg("InMemory auth repo")
	return &repo{
		mx: &sync.Mutex{},
		db: map[string]*ddd.AuthPasswd{},
		h:  make([]ddd.AuthHistory, 0),
	}, nil
}

func (r *repo) GetPasswd(ctx context.Context, email string) ([]ddd.AuthPasswd, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	a, ok := r.db[email]
	if !ok {
		return []ddd.AuthPasswd{}, nil
	}
	return []ddd.AuthPasswd{*a}, nil
}

func (r *repo) ListPasswd(ctx context.Context, filter app.FilterFactory) ([]ddd.AuthPasswd, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	out := []ddd.AuthPasswd{}
	for _, v := range r.db {
		out = append(out, *v)
	}
	return out, nil
}

func (r *repo) SavePasswd(ctx context.Context, auth ddd.AuthPasswd) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.db[auth.Subject] = &auth
	logger.Warn().Str("subject", auth.Subject).Msg("Save")
	return nil
}

func (r *repo) Update(ctx context.Context, auth ddd.AuthPasswd) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	old, ok := r.db[auth.Subject]
	if !ok {
		return nil
	}
	old.Status = auth.Status
	old.SystemRoles = auth.SystemRoles

	// r.db[auth.Email] = old
	logger.Warn().Str("subject", auth.Subject).Msg("Save")
	return nil
}

func (r *repo) UpdatePassword(ctx context.Context, email string, password string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	v, ok := r.db[email]
	if !ok {
		return fmt.Errorf("subject not found")
	}

	v.Passwd = password
	return nil

}

func (r *repo) UpdateStatus(ctx context.Context, email string, status ddd.StatusType) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	v, ok := r.db[email]
	if !ok {
		return fmt.Errorf("email not found")
	}

	v.Status = status
	return nil
}
