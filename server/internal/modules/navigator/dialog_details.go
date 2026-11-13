package navigator

import "github.com/babakapaa/lodestar/server/internal/core/module"

type DialogDetails struct {
	Dialog
	Brief    *module.Brief `json:"brief"`
	Messages []Message     `json:"messages"`
}
