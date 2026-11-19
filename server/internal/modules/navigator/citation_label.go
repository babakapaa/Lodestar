package navigator

func (c Citation) label() string {
	if c.Reference != "" {
		return c.Reference
	}
	return c.Title
}
