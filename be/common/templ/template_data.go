package templ

type (
	TableData struct {
		Name    string
		Headers []string
		Rows    [][]string
	}

	TemplateData struct {
		Person map[string]any
		Items  map[string]any
		Lists  map[string][]string
		Tables map[string]TableData
	}

	DataConverter interface {
		ConvertToKey(itemKey string, dest *TemplateData)
	}
)
