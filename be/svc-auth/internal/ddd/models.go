package ddd

import (
	"context"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	AuthToken struct {
		Token     string
		ExpiresAt time.Time
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
)

const (
	StatusEnabled StatusType = "enabled"
	StatusDisable StatusType = "disable"
	StatusPending StatusType = "pending"
)

type AuthPasswdRepo interface {
	Save(ctx context.Context, auth AuthPasswd) error
	Get(ctx context.Context, email string) ([]AuthPasswd, error)
	List(ctx context.Context, filter app.FilterFactory) ([]AuthPasswd, error)
	UpdatePassword(ctx context.Context, email, password string) error
	UpdateStatus(ctx context.Context, email string, status StatusType) error
	Update(ctx context.Context, auth AuthPasswd) error
}

func (a *AuthPasswd) ToAuthInfo(domain string) roles.AuthInfo {
	out := roles.AuthInfo{
		User:        a.Email,
		Realm:       domain,
		Roles:       []string{},
		SystemRoles: []string{},
	}

	hostRoles, ok := a.RealmRoles[domain]
	if ok {
		out.Roles = append(out.Roles, hostRoles...)
	}
	out.SystemRoles = append(out.SystemRoles, a.SystemRoles...)

	return out
}
