package ddd

import (
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
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

var _ (logger.TraceDataExtractor) = (Template)(Template{})

func (t Template) Extract() logger.TraceData {
	out := logger.TraceData{
		"tmpl.Name":        logger.TraceValueString(t.Name),
		"tmpl.ContentType": logger.TraceValueString(t.ContentType),
	}
	if t.Id == "" {
		out["tmpl.id"] = logger.TraceValueString(t.Id)
	}
	return out
}
