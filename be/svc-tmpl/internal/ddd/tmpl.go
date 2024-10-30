package ddd

import (
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type TemplateTable struct {
	Name    string
	Headers []string
	Rows    [][]string
}

type TemplateData struct {
	UserInfo roles.AuthInfo
	Items    map[string]any
	Lists    map[string][]string
	Tables   map[string]TemplateTable
}

type Template struct {
	Id          string
	ContentType string
	Name        string
	Labels      []string
	Attachments map[string]string
	Images      map[string]string
	Body        string
	CreatedAt   time.Time
}

func (t Template) Extractor() map[string]string {
	out := map[string]string{
		"tmpl.Name":        t.Name,
		"tmpl.ContentType": t.ContentType,
	}
	if t.Id == "" {
		out["tmpl.id"] = t.Id
	}
	return out
}
