package extauth

import (
	"context"
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"golang.org/x/oauth2"
)

type Auth2Profile struct {
	ID      string
	Email   string
	Picture string
	Attr    map[string]string
}

type LoginRequest struct {
	State       string
	Code        string
	Provider    string
	RedirectUrl string
}

type ProviderConfig struct {
	ClientID        string
	ClientSecret    string
	Scopes          []string
	RedirectURL     string
	AuthURL         string
	TokenURL        string
	AuthStyle       string //: oauth2.AuthStyleInParams,
	FetchProfileURL string
}

type auth2 struct {
	cfg ProviderConfig
}

func New(cfg ProviderConfig) *auth2 {
	return &auth2{
		cfg: cfg,
	}
}

func (a *auth2) CodeExchange(ctx context.Context, req LoginRequest) (*http.Client, error) {
	conf := oauth2.Config{
		ClientID:     a.cfg.ClientID,
		ClientSecret: a.cfg.ClientSecret,
		Scopes:       a.cfg.Scopes,
		// RedirectURL:  "http://localhost:8080./callback",
		RedirectURL: req.RedirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL: a.cfg.AuthURL,
			// TokenURL: "https://oauth2.googleapis.com/token",
			// TokenURL: "https://www.googleapis.com/oauth2/v4/token:",
			TokenURL:  a.cfg.TokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	token, err := conf.Exchange(ctx, req.Code)
	if err != nil {
		log.Log().Error(err, "conf.Exchange")
		return nil, err
	}

	return conf.Client(ctx, token), nil
}
