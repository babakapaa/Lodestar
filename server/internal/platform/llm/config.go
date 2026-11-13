package llm

import "time"

type Config struct {
	URL     string
	Model   string
	Key     string
	Timeout time.Duration
}

func (c Config) Enabled() bool {
	return c.URL != "" && c.Model != ""
}
