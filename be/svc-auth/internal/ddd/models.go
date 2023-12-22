package ddd

import (
	"context"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	StatusType string

	AuthPasswd struct {
		Email       string     `json:"email"`
		Passwd      string     `json:"passwd"`
		Status      StatusType `json:"status"`
		SystemRoles []string   `json:"system_roles"`
		CreatedAt   time.Time  `json:"created_at"`
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
	List(ctx context.Context) ([]AuthPasswd, error)
	UpdatePassword(ctx context.Context, email, password string) error
	UpdateStatus(ctx context.Context, email string, status StatusType) error
	Update(ctx context.Context, auth AuthPasswd) error
}

func (a *AuthPasswd) ToAuthInfo() roles.AuthInfo {
	out := roles.AuthInfo{
		User:  a.Email,
		Roles: []roles.RoleName{},
	}
	for _, v := range a.SystemRoles {
		out.Roles = append(out.Roles, roles.RoleName(v))
	}
	return out
}
