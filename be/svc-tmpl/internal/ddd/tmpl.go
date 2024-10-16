package ddd

import (
	"io"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type FileWriterTo struct {
	ContentType string
	Version     string
	Name        string
	WriterTo    io.WriterTo
}

type DataTable struct {
	Name    string
	Headers []string
	Rows    [][]string
}

type TemplateData struct {
	UserInfo roles.AuthInfo
	Items    map[string]any
	Lists    map[string][]string
	Tables   map[string]DataTable
}

type Template struct {
	ContentType string
	Version     string
	Name        string
	Body        string
	Attachments map[string]string
	Images      map[string]string
}
