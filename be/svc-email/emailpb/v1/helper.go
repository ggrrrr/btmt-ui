package emailpbv1

import (
	"log/slog"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

func CreateList(l []*EmailAddr) []string {
	if len(l) == 0 {
		return []string{}
	}

	list := make([]string, 0, len(l))

	for _, v := range l {
		list = append(list, v.Email)
	}
	return list

}

var _ (td.TraceDataExtractor) = (*SendEmail)(nil)

// Extract implements logger.TraceDataExtractor.
func (x *SendEmail) Extract() *td.TraceData {
	return &td.TraceData{
		KV: map[string]slog.Value{
			"email.id":            slog.StringValue(x.Id),
			"email.account.realm": slog.StringValue(x.Message.FromAccount.Realm),
		},
	}
}
