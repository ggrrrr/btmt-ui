package tmplpbv1

import (
	"log/slog"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

var _ (td.TraceDataExtractor) = (*Template)(nil)

func (t *Template) Extract() *td.TraceData {
	out := &td.TraceData{
		KV: map[string]slog.Value{
			"tmpl.Name":        slog.StringValue(t.Name),
			"tmpl.ContentType": slog.StringValue(t.ContentType),
		},
	}
	if t.Id != "" {
		out.KV["tmpl.Id"] = slog.StringValue(t.Id)
	}
	return out
}

var _ (td.TraceDataExtractor) = (*TemplateUpdate)(nil)

func (t *TemplateUpdate) Extract() *td.TraceData {
	out := &td.TraceData{
		KV: map[string]slog.Value{
			"tmpl.Name":        slog.StringValue(t.Name),
			"tmpl.ContentType": slog.StringValue(t.ContentType),
		}}
	if t.Id != "" {
		out.KV["tmpl.Id"] = slog.StringValue(t.Id)
	}
	return out
}
