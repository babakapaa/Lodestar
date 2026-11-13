package agent

import (
	"errors"
	"time"
)

type Config struct {
	Server  string
	Token   string
	GPUs    int
	Host    string
	Scheme  string
	DataDir string
	Every   time.Duration
}

func (c Config) Validate() error {
	switch {
	case c.Server == "":
		return errors.New("LODESTAR_SERVER is required")
	case c.Token == "":
		return errors.New("LODESTAR_NODE_TOKEN is required")
	case c.GPUs < 1:
		return errors.New("LODESTAR_AGENT_GPUS must be at least 1")
	case c.Host == "":
		return errors.New("LODESTAR_AGENT_HOST is required: the address students open sessions at")
	case c.DataDir == "":
		return errors.New("LODESTAR_AGENT_DATA is required")
	}
	return nil
}
