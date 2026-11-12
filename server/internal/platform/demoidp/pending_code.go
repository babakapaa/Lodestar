package demoidp

type pendingCode struct {
	redirectURI string
	nonce       string
	challenge   string
	user        User
}
