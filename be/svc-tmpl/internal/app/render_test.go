package app

import (
	"context"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

func TestRender(t *testing.T) {

	ctx := roles.CtxWithAuthInfo(context.Background(), roles.AuthInfo{
		Subject: "user@me",
		Realm:   "localshit",
	})

	testapp := &App{}

	testapp.RenderHtml(ctx, RenderRequest{
		Items: map[string]any{
			"someData": "some value",
		},
		TemplateBody: `me: {{ .UserInfo.User }}`,
	})
}
