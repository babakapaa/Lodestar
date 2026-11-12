package auth

type idClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
