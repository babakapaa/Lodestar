package auth

type LoginAttempt struct {
	State         string
	ProviderID    string
	Nonce         string
	Verifier      string
	RedirectURI   string
	CodeChallenge string
	DeviceName    string
}
