package templ

import (
	"fmt"
	htmltemplate "html/template"
	"io"

	templv1 "github.com/ggrrrr/btmt-ui/be/common/templ/v1"
)

type OptionFunc func(*HtmlTemplate) error

type (
	HtmlTemplate struct {
		htmlTempl *htmltemplate.Template
	}
)

const (
	templateRenderImg string = "renderImg"
)

func NewHtml(body string, opts ...OptionFunc) (*HtmlTemplate, error) {
	var err error
	out := &HtmlTemplate{
		htmlTempl: htmltemplate.New("base"),
	}

	for _, optfunc := range opts {
		err = optfunc(out)
		if err != nil {
			return nil, err
		}
	}

	out.htmlTempl, err = out.htmlTempl.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse body %w", err)
	}

	return out, nil
}

func (t *HtmlTemplate) Execute(writer io.Writer, data *templv1.Data) error {
	return t.htmlTempl.Execute(writer, fromV1(data))
}

func WithRenderImg(callF func(imageName string) htmltemplate.HTML) OptionFunc {
	return func(t *HtmlTemplate) error {
		if t == nil {
			return fmt.Errorf("nil templ")
		}
		if t.htmlTempl != nil {
			t.htmlTempl = t.htmlTempl.Funcs(htmltemplate.FuncMap{
				templateRenderImg: callF,
			})
		}
		return nil
	}
}
