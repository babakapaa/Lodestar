package livekit

type Config struct {
	URL       string
	APIKey    string
	APISecret string
}

func (c Config) Enabled() bool {
	return c.URL != "" && c.APIKey != "" && c.APISecret != ""
}
