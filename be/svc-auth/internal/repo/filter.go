package repo

import "github.com/ggrrrr/btmt-ui/be/common/app"

type AuthFilter struct {
	includeTenants []string
}

var _ (app.FilterFactory) = (*AuthFilter)(nil)

func (a *AuthFilter) Create() any {
	return a.includeTenants
}

func NewFilter(tenants ...string) (*AuthFilter, error) {
	out := &AuthFilter{
		includeTenants: tenants,
	}

	return out, nil
}
