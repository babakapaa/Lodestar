package chat

func callLimit(c Chat) int {
	if c.Kind == KindChannel {
		return MaxChannelCallPeople
	}
	return MaxCallPeople
}
