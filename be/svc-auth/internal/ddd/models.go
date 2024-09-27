package ddd

import (
	"context"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	StatusType string

	AuthPasswd struct {
		Email       string              `json:"email"`
		Passwd      string              `json:"passwd"`
		Status      StatusType          `json:"status"`
		TenantRoles map[string][]string `json:"tenant_roles"`
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

func (a *AuthPasswd) ToAuthInfo(tenant roles.Tenant) roles.AuthInfo {
	out := roles.AuthInfo{
		User:        a.Email,
		Tenant:      tenant,
		Roles:       []roles.RoleName{},
		SystemRoles: []roles.RoleName{},
	}
	hostRoles := a.TenantRoles[string(tenant)]
	for _, v := range hostRoles {
		out.Roles = append(out.Roles, roles.RoleName(v))
	}
	for _, v := range a.SystemRoles {
		out.SystemRoles = append(out.SystemRoles, roles.RoleName(v))
	}
	return out
}
