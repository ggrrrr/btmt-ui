package app

import (
	"bytes"
	"context"
	"fmt"
	htmltemplate "html/template"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func (a *App) RenderHtml(ctx context.Context, template ddd.Template, data ddd.TemplateData) (string, error) {

	authInfo := roles.AuthInfoFromCtx(ctx)

	tmplValidator := validator(ctx, authInfo.Realm, a)
	tmplValidator.resized = true

	tmpl, err := htmltemplate.New("template_data").
		Funcs(htmltemplate.FuncMap{
			"renderImg": tmplValidator.RenderImg,
		}).
		Parse(template.Body)
	if err != nil {
		return "", err
	}

	if len(tmplValidator.errors) > 0 {
		fmt.Printf("\n\n, %+v \n", tmplValidator.errors)
	}

	buf := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(buf, data)

	return buf.String(), err

}
