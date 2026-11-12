package auth

import "time"

type Lifetimes struct {
	Access  time.Duration
	Refresh time.Duration
	Login   time.Duration
	Grant   time.Duration
}
