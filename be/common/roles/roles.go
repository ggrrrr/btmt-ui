package roles

const AuthSchemeBearer string = "Bearer"

type (
	ctxKeyType struct{}

	RoleName string

	Tenant string

	AuthCredentials string

	Authorization struct {
		AuthScheme      string          // Basic, Bearer, ...,
		AuthCredentials AuthCredentials // JWT TOKEN OR OTHER secret data
	}
	UserRequest struct {
		FullMethod    string
		Device        Device
		Authorization Authorization
	}
)

const (
	RoleAdmin    RoleName = "admin"
	SystemTenant Tenant   = "localhost"
)

func HasRole(role RoleName, roles []RoleName) bool {
	for r := range roles {
		if roles[r] == role {
			return true
		}
	}
	return false
}
