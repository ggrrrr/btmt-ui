package ddd

import (
	"context"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/google/uuid"
)

type (
	AuthToken struct {
		Value     string
		ExpiresAt time.Time
	}

	LoginToken struct {
		ID           uuid.UUID
		Email        string
		AccessToken  AuthToken
		RefreshToken AuthToken
	}

	StatusType string

	AuthPasswd struct {
		Email       string              `json:"email"`
		Passwd      string              `json:"passwd"`
		Status      StatusType          `json:"status"`
		RealmRoles  map[string][]string `json:"realm_roles"`
		SystemRoles []string            `json:"system_roles"`
		CreatedAt   time.Time           `json:"created_at"`
	}

	AuthHistory struct {
		ID        uuid.UUID
		User      string
		Method    string
		Device    roles.Device
		CreatedAt time.Time
	}
)

const (
	StatusEnabled StatusType = "enabled"
	StatusDisable StatusType = "disable"
	StatusPending StatusType = "pending"
)

type AuthRepo interface {
	AuthPasswdRepo
	AuthHistoryRepo
}

type AuthPasswdRepo interface {
	SavePasswd(ctx context.Context, auth AuthPasswd) error
	GetPasswd(ctx context.Context, email string) ([]AuthPasswd, error)
	ListPasswd(ctx context.Context, filter app.FilterFactory) ([]AuthPasswd, error)
	UpdatePassword(ctx context.Context, email, password string) error
	UpdateStatus(ctx context.Context, email string, status StatusType) error
	Update(ctx context.Context, auth AuthPasswd) error
}

type AuthHistoryRepo interface {
	SaveHistory(ctx context.Context, info roles.AuthInfo, method string) (err error)
	ListHistory(ctx context.Context, user string) (authHistory []AuthHistory, err error)
	GetHistory(ctx context.Context, id uuid.UUID) (authHistory *AuthHistory, err error)
	DeleteHistory(ctx context.Context, id string) (err error)
}

func (a *AuthPasswd) ToAuthInfo(device roles.Device, domain string) roles.AuthInfo {
	out := roles.AuthInfo{
		User:        a.Email,
		Realm:       domain,
		Device:      device,
		Roles:       []string{},
		SystemRoles: []string{},
		ID:          uuid.New(),
	}

	hostRoles, ok := a.RealmRoles[domain]
	if ok {
		out.Roles = append(out.Roles, hostRoles...)
	}
	out.SystemRoles = append(out.SystemRoles, a.SystemRoles...)

	return out
}
