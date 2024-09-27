package roles

const AuthSchemeBearer string = "Bearer"

type (
	ctxKeyType struct{}

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
	RoleAdmin    string = "admin"
	SystemTenant string = "localhost"
)

func HasRole(role string, roles []string) bool {
	for r := range roles {
		if roles[r] == role {
			return true
		}
	}
	return false
}
