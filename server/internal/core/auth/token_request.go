package auth

type TokenRequest struct {
	GrantType    string `json:"grantType"`
	Code         string `json:"code"`
	CodeVerifier string `json:"codeVerifier"`
	RedirectURI  string `json:"redirectUri"`
	RefreshToken string `json:"refreshToken"`
}
