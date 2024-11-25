package ddd

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
