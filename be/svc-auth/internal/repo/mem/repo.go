package mem

import (
	"context"
	"fmt"
	"sync"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

type (
	repo struct {
		mx *sync.Mutex
		db map[string]*ddd.AuthPasswd
	}
)

var _ (ddd.AuthPasswdRepo) = (*repo)(nil)

func New() *repo {
	logger.Warn().Msg("InMemmory auth repo")
	return &repo{
		mx: &sync.Mutex{},
		db: map[string]*ddd.AuthPasswd{},
	}
}

func (r *repo) Get(ctx context.Context, email string) ([]ddd.AuthPasswd, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	a, ok := r.db[email]
	if !ok {
		return []ddd.AuthPasswd{}, nil
	}
	return []ddd.AuthPasswd{*a}, nil
}

func (r *repo) List(ctx context.Context) ([]ddd.AuthPasswd, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	out := []ddd.AuthPasswd{}
	for _, v := range r.db {
		out = append(out, *v)
	}
	return out, nil
}

func (r *repo) Save(ctx context.Context, auth ddd.AuthPasswd) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.db[auth.Email] = &auth
	logger.Warn().Str("email", auth.Email).Msg("Save")
	return nil
}

func (r *repo) UpdatePassword(ctx context.Context, email string, password string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	v, ok := r.db[email]
	if !ok {
		return fmt.Errorf("email not found")
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