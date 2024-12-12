package emailpbv1

import "github.com/ggrrrr/btmt-ui/be/common/logger"

func (c *EmailMessage) CreateToList() []string {
	if len(c.ToEmail) == 0 {
		return []string{}
	}

	list := make([]string, 0, len(c.ToEmail))

	for _, v := range c.ToEmail {
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
