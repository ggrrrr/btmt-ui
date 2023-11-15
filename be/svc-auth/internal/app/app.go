package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/stackus/errors"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthToken string

	AppCfgFunc func(a *application) error

	App interface {
		LoginPasswd(ctx context.Context, email, passwd string) (app.Result[AuthToken], error)

		UpdatePasswd(ctx context.Context, email, oldPasswd, newPasswd string) error
		Validate(ctx context.Context) error
		ListAuth(ctx context.Context) (app.Result[[]ddd.AuthPasswd], error)
		CreateAuth(ctx context.Context, auth ddd.AuthPasswd) error
		// RegisterEmail(ctx context.Context, email string) (*a.Result[string], error)
		// EnableEmail(ctx context.Context, email string) (*a.Result[string], error)
		// DisableAuth(ctx context.Context, email string) (a.Result[string], error)
	}

	application struct {
		appPolices roles.AppPolices
		authRepo   ddd.AuthPasswdRepo
		signer     token.Signer
	}
)

var _ (App) = (*application)(nil)

func New(cfgs ...AppCfgFunc) (*application, error) {
	a := &application{}
	for _, c := range cfgs {
		err := c(a)
		if err != nil {
			return nil, err
		}
	}
	if a.appPolices == nil {
		logger.Log().Warn().Msg("use mock AppPolices")
		a.appPolices = roles.NewAppPolices()
	}
	return a, nil
}

func WithAuthRepo(repo ddd.AuthPasswdRepo) AppCfgFunc {
	return func(a *application) error {
		a.authRepo = repo
		return nil
	}
}

func WithTokenSigner(s token.Signer) AppCfgFunc {
	return func(a *application) error {
		a.signer = s
		return nil
	}
}

func canLogin(auth *ddd.AuthPasswd) bool {
	if auth == nil {
		return false
	}
	if auth.Status == ddd.StatusEnabled {
		return true
	}
	return false
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (ap *application) findEmail(ctx context.Context, email string) (*ddd.AuthPasswd, error) {
	auths, err := ap.authRepo.Get(ctx, email)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternalServerError, err.Error())
	}
	if len(auths) == 0 {
		return nil, nil
	}
	if len(auths) > 1 {
		logger.Log().Error().Str("email", string(email)).Msg("multiple result")
		return nil, ErrAuthMultipleEmail
	}
	return &auths[0], nil
}
