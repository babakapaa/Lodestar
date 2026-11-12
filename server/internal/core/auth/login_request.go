package auth

type LoginRequest struct {
	TenantSlug    string
	ProviderKey   string
	RedirectURI   string
	CodeChallenge string
	DeviceName    string
}
