package auth

type LoginGrant struct {
	TenantID      string
	UserID        string
	RedirectURI   string
	CodeChallenge string
	DeviceName    string
}
