package auth

type Callback struct {
	TenantSlug  string
	ProviderKey string
	State       string
	Code        string
	Error       string
}
