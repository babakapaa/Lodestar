package mail

type Config struct {
	Addr     string
	Username string
	Password string
	From     string
}

func (c Config) Enabled() bool {
	return c.Addr != "" && c.From != ""
}
