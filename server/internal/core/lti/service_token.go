package lti

import "time"

type serviceToken struct {
	value   string
	expires time.Time
}
