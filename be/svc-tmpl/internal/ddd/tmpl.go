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
	Id          string            `json:"id,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Name        string            `json:"name,omitempty"`
	Labels      []string          `json:"labels,omitempty"`
	Attachments map[string]string `json:"attachments,omitempty"`
	Images      map[string]string `json:"images,omitempty"`
	Body        string            `json:"body,omitempty"`
	CreatedAt   time.Time         `json:"created_at,omitempty"`
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
