package tmplpbv1

import "github.com/ggrrrr/btmt-ui/be/common/logger"

var _ (logger.TraceDataExtractor) = (*Template)(nil)

func (t *Template) Extract() logger.TraceData {
	out := logger.TraceData{
		"tmpl.Name":        logger.TraceValueString(t.Name),
		"tmpl.ContentType": logger.TraceValueString(t.ContentType),
	}
	if t.Id != "" {
		out["tmpl.Id"] = logger.TraceValueString(t.Id)
	}
	return out
}

var _ (logger.TraceDataExtractor) = (*TemplateUpdate)(nil)

func (t *TemplateUpdate) Extract() logger.TraceData {
	out := logger.TraceData{
		"tmpl.Name":        logger.TraceValueString(t.Name),
		"tmpl.ContentType": logger.TraceValueString(t.ContentType),
	}
	if t.Id != "" {
		out["tmpl.Id"] = logger.TraceValueString(t.Id)
	}
	return out
}
