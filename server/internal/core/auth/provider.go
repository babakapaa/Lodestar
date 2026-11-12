package auth

type Provider struct {
	ID           string
	TenantID     string
	Key          string
	DisplayName  string
	Issuer       string
	ClientID     string
	ClientSecret string
	Scopes       []string
}
