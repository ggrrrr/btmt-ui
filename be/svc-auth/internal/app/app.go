package app

import (
	"context"
	"fmt"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/stackus/errors"
	"golang.org/x/crypto/bcrypt"
)

type (
	AppCfgFunc func(a *Application) error

	App interface {
		Get(ctx context.Context, email string) (ddd.AuthPasswd, error)
		UserCreate(ctx context.Context, auth ddd.AuthPasswd) error
		UserList(ctx context.Context) ([]ddd.AuthPasswd, error)
		UserUpdate(ctx context.Context, email ddd.AuthPasswd) error
		UserChangePasswd(ctx context.Context, email, oldPasswd, newPasswd string) error

		LoginPasswd(ctx context.Context, email, passwd string) (ddd.LoginToken, error)
		TokenValidate(ctx context.Context) error
		TokenRefresh(ctx context.Context) (ddd.LoginToken, error)

		// RegisterEmail(ctx context.Context, email string) (*a.Result[string], error)
		// EnableEmail(ctx context.Context, email string) (*a.Result[string], error)
		// DisableAuth(ctx context.Context, email string) (a.Result[string], error)
	}

	Application struct {
		accessTokenTTL  time.Duration
		refreshTokenTTL time.Duration
		appPolices      roles.AppPolices
		authRepo        ddd.AuthPasswdRepo
		historyRepo     ddd.AuthHistoryRepo
		signer          token.Signer
	}
)

var _ (App) = (*Application)(nil)

func New(cfgs ...AppCfgFunc) (*Application, error) {
	a := &Application{}
	for _, c := range cfgs {
		err := c(a)
		if err != nil {
			return nil, err
		}
	}
	if a.appPolices == nil {
		logger.Warn().Msg("use mock AppPolices")
		a.appPolices = roles.NewAppPolices()
	}
	logger.Info().
		Int("ttl.refresh.days", int(a.refreshTokenTTL.Hours()/24)).
		Int("ttl.token.minutes", int(a.accessTokenTTL.Minutes())).
		Msg("app.New")
	return a, nil
}

func WithAuthRepo(repo ddd.AuthPasswdRepo) AppCfgFunc {
	return func(a *Application) error {
		if repo == nil {
			return fmt.Errorf("repo is nil")
		}
		a.authRepo = repo
		return nil
	}
}

func WithHistoryRepo(repo ddd.AuthHistoryRepo) AppCfgFunc {
	return func(a *Application) error {
		if repo == nil {
			return fmt.Errorf("repo is nil")
		}
		a.historyRepo = repo
		return nil
	}
}

func WithTokenSigner(s token.Signer) AppCfgFunc {
	return func(a *Application) error {
		a.signer = s
		return nil
	}
}

func WithTokenTTL(tokenTTL time.Duration, refreshTTL time.Duration) AppCfgFunc {
	return func(a *Application) error {
		if tokenTTL == 0 {
			return fmt.Errorf("ttl.token is 0")
		}
		if refreshTTL == 0 {
			return fmt.Errorf("ttl.refresh is 0")
		}
		a.accessTokenTTL = tokenTTL
		a.refreshTokenTTL = refreshTTL
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (ap *Application) findEmail(ctx context.Context, email string) (*ddd.AuthPasswd, error) {
	auths, err := ap.authRepo.GetPasswd(ctx, email)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternalServerError, err.Error())
	}
	if len(auths) == 0 {
		return nil, nil
	}
	if len(auths) > 1 {
		logger.Error(fmt.Errorf("multiple result")).Str("email", string(email)).Msg("findEmail")
		return nil, ErrAuthMultipleEmail
	}
	return &auths[0], nil
}
