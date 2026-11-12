package demoidp

type Config struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	Redirects    []string
	Users        []User
	Pick         func() User
}
