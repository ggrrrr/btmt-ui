package ddd

import (
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type TemplateTable struct {
	Name    string
	Headers []string
	Rows    [][]string
}

type TemplateData struct {
	Person map[string]string
	Items  map[string]any
	Lists  map[string][]string
	Tables map[string]TemplateTable
}

type Template struct {
	Id          string            `json:"id,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Name        string            `json:"name,omitempty"`
	Labels      []string          `json:"labels,omitempty"`
	Images      []string          `json:"images,omitempty"`
	Files       map[string]string `json:"files,omitempty"`
	Body        string            `json:"body,omitempty"`
	CreatedAt   time.Time         `json:"created_at,omitempty"`
	UpdatedAt   time.Time         `json:"updated_at,omitempty"`
}

// We need this implemented coz of easy tracing
var _ (logger.TraceDataExtractor) = (Template)(Template{})

func (t Template) Extract() logger.TraceData {
	out := logger.TraceData{
		"tmpl.Name":        logger.TraceValueString(t.Name),
		"tmpl.ContentType": logger.TraceValueString(t.ContentType),
	}
	if t.Id != "" {
		out["tmpl.Id"] = logger.TraceValueString(t.Id)
	}
	return out
}
