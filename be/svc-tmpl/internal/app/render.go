package app

import (
	"bytes"
	"context"
	"fmt"
	htmltemplate "html/template"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

type RenderRequest struct {
	Items        map[string]any
	TemplateBody string
}

func (a *App) RenderHtml(ctx context.Context, request RenderRequest) (*ddd.FileWriterTo, error) {

	authInfo := roles.AuthInfoFromCtx(ctx)

	data := ddd.TemplateData{
		UserInfo: authInfo,
		Items:    request.Items,
	}

	tmpl, err := htmltemplate.New("template_data").Parse(request.TemplateBody)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(buf, data)
	// tmpl := template.Mu

	fmt.Printf("result %s\n", buf.String())
	return nil, err

}
