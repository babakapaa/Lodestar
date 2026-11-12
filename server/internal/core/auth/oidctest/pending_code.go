package oidctest

type pendingCode struct {
	clientID    string
	redirectURI string
	nonce       string
	challenge   string
	user        User
}
