package roles

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

const AuthSchemeBearer string = "Bearer"

type (
	authInfoCtxKeyType struct{}

	AuthInfo struct {
		// User name or system name
		Subject string
		// For `sudo` like behavior
		AdminSubject string
		// Domain name
		Realm string
		// List of domain roles
		Roles []string
		// List of system roles
		SystemRoles []string
		// Info of the agent device ( browser,app,service,etc...)
		Device app.Device
		// Unique ID of the token
		ID uuid.UUID
	}
)

var _ (td.TraceDataExtractor) = AuthInfo{}

// Extract implements td.TraceDataExtractor.
func (a AuthInfo) Extract() *td.TraceData {
	if a.Subject == "" {
		return nil
	}
	return &td.TraceData{
		KV: map[string]slog.Value{
			"auth.info": slog.GroupValue(
				slog.String("subject", a.Subject),
				slog.String("realm", a.Realm),
				slog.String("device.info", a.Device.DeviceInfo),
				slog.String("device.remote.addr", a.Device.RemoteAddr),
			),
		},
	}
}

func CreateSystemAdminUser(tenant string, subject string, device app.Device) AuthInfo {
	return AuthInfo{
		Realm:       tenant,
		Subject:     subject,
		Roles:       []string{RoleAdmin},
		SystemRoles: []string{RoleAdmin},
		Device:      device,
	}
}

func CtxWithAuthInfo(ctx context.Context, authInfo AuthInfo) context.Context {
	return context.WithValue(ctx, authInfoCtxKeyType{}, authInfo)
}

func AuthInfoFromCtx(ctx context.Context) AuthInfo {
	value, ok := ctx.Value(authInfoCtxKeyType{}).(AuthInfo)
	if !ok {
		return AuthInfo{}
	}
	return value
}
