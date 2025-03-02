package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.be.svc-auth"

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
		otelTracer      tracer.OTelTracer
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
	a := &Application{
		otelTracer: tracer.Tracer(otelScope),
	}
	for _, c := range cfgs {
		err := c(a)
		if err != nil {
			return nil, err
		}
	}
	if a.appPolices == nil {
		log.Log().Warn(nil, "use mock AppPolices")
		a.appPolices = roles.NewAppPolices()
	}
	if a.accessTokenTTL == 0 {
		return nil, fmt.Errorf("ttl.token is 0")
	}
	if a.refreshTokenTTL == 0 {
		return nil, fmt.Errorf("ttl.refresh is 0")
	}

	log.Log().Info("app.New",
		slog.Int("refresh.token.ttl.days", int(a.refreshTokenTTL.Hours()/24)),
		slog.Int("access.token.ttl.minutes", int(a.accessTokenTTL.Minutes())))
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

func WithTokenTTL(accessTokenTTL time.Duration, refreshTokenTTL time.Duration) AppCfgFunc {
	return func(a *Application) error {
		a.accessTokenTTL = accessTokenTTL
		a.refreshTokenTTL = refreshTokenTTL
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
		return nil, app.SystemError("please try again later", err)
	}
	if len(auths) == 0 {
		return nil, nil
	}
	if len(auths) > 1 {
		log.Log().ErrorCtx(ctx, fmt.Errorf("multiple result"), "findEmail")
		return nil, ErrAuthMultipleEmail
	}
	return &auths[0], nil
}
