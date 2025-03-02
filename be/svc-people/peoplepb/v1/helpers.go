package peoplepbv1

import (
	"log/slog"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

// We need this implemented coz of easy tracing
var _ (td.TraceDataExtractor) = (*Person)(nil)

func (p *Person) Extract() *td.TraceData {
	return &td.TraceData{
		KV: map[string]slog.Value{
			"person.id": slog.StringValue(p.Id),
		},
	}
}

func (f *ListRequest) ToFilter() map[string][]string {
	out := map[string][]string{}
	if f.Filters == nil {
		return out
	}
	for k, v := range f.Filters {
		if v != nil {
			out[k] = v.GetList()
		}
	}
	return out
}
