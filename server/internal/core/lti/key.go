package lti

import "crypto/rsa"

type Key struct {
	ID      string
	Private *rsa.PrivateKey
}
