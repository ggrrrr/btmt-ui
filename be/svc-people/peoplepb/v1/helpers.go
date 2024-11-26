package peoplepbv1

import (
	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

// We need this implemented coz of easy tracing
var _ (logger.TraceDataExtractor) = (*Person)(nil)

func (p Person) Extract() logger.TraceData {
	return logger.TraceData{
		"person.id": logger.TraceValueString(p.Id),
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
