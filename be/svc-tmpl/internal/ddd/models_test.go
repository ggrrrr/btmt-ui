package ddd

import (
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/stretchr/testify/require"
)

func TestData(t *testing.T) {

	testTmpl := Tmpl{
		ContentType: "text/plain",
		Body: `# Header from: {{ .UserInfo.User }}
# Header key 1: {{ .Items.key1.Item1 }} {{ .Items.key1.Item }}
# List {{range index .Lists "list1"}}
	* {{.}}{{end}}
---
# table {{ .Tables.table1.Name }}
Header: {{range .Tables.table1.Headers}} {{ .}} {{end}}
{{range .Tables.table1.Rows}}row   : {{ range .}} {{ . }}{{ end}} 
{{end}}end
`,
	}

	item1 := struct {
		Item1 string
		Item  string
	}{Item1: "value 1", Item: "value 2"}
	item2 := struct{ Item2 string }{Item2: "value 2"}

	testData := TmplData{
		UserInfo: roles.AuthInfo{
			User: "test user",
		},
		Items: map[string]any{
			"key1": item1,
			"key2": item2,
		},
		Lists: map[string][]string{"list1": {"item 1", "item 2", "item 3"}},
		Tables: map[string]TmplTable{
			"table1": {
				Name:    "table 1",
				Headers: []string{"name 1", "name 2"},
				Rows: [][]string{
					{"cel 11", "cel 12"},
					{"cel 21", "cel 22"},
				},
			},
		},
	}

	tmpl := template.Must(template.New("template_data").Parse(testTmpl.Body))
	require.NotNil(t, tmpl)

	buffer := new(strings.Builder)
	tmpl.Execute(buffer, testData)
	fmt.Printf("----\n%s\n-----\n", buffer.String())
}
