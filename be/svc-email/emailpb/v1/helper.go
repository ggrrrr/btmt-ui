package emailpbv1

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
