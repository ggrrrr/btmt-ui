package ddd

import (
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	polices := roles.NewAppPolices()

	tests := []struct {
		name     string
		domain   string
		authPass AuthPasswd
		expErr   error
	}{
		{
			name:   "system admin",
			domain: "domain",
			authPass: AuthPasswd{
				Email:       "asdasd",
				Passwd:      "",
				Status:      "",
				TenantRoles: map[string][]string{},
				SystemRoles: []string{"asdasd", roles.RoleAdmin},
			},
			expErr: nil,
		},
		{
			name:   "domain admin role ok",
			domain: "domain",
			authPass: AuthPasswd{
				Email:  "asdasd",
				Passwd: "",
				Status: "",
				TenantRoles: map[string][]string{
					"domain": {roles.RoleAdmin, "asdasd"},
				},
				SystemRoles: []string{"asdasd", "asd"},
			},
			expErr: nil,
		},
		{
			name:   "domain not match no role",
			domain: "domain",
			authPass: AuthPasswd{
				Email:  "asdasd",
				Passwd: "",
				Status: "",
				TenantRoles: map[string][]string{
					"bad-domain": {roles.RoleAdmin, "asdasd"},
				},
				SystemRoles: []string{"asdasd", "asd"},
			},
			expErr: app.ErrForbidden,
		},
		{
			name:   "no roles",
			domain: "domain",
			authPass: AuthPasswd{
				Email:       "asdasd",
				Passwd:      "",
				Status:      "",
				TenantRoles: map[string][]string{},
				SystemRoles: []string{"asdasd"},
			},
			expErr: app.ErrForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authInfo := tc.authPass.ToAuthInfo(tc.domain)
			err := polices.CanDo(tc.domain, "somepath", authInfo)
			if tc.expErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expErr)
			} else {
				assert.NoError(t, err)
			}
		})

	}

}
