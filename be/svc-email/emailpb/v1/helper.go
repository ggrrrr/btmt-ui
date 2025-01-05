package emailpbv1

import "github.com/ggrrrr/btmt-ui/be/common/logger"

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

var _ (logger.TraceDataExtractor) = (*SendEmail)(nil)

// Extract implements logger.TraceDataExtractor.
func (x *SendEmail) Extract() logger.TraceData {
	return logger.TraceData{
		"email.id":            logger.TraceValueString(x.Id),
		"email.account.realm": logger.TraceValueString(x.Message.FromAccount.Realm),
	}
}
