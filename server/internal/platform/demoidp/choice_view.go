package demoidp

import "net/url"

type choiceView struct {
	Users []User
	Query url.Values
}
