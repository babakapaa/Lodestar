package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
)

func coreSettings() []schema.Definition {
	var list []schema.Definition
	for _, part := range [][]schema.Definition{
		settings.CoreSettings(), files.CoreSettings(), delivery.CoreSettings(), audit.CoreSettings(), privacy.CoreSettings(),
	} {
		list = append(list, part...)
	}
	return list
}
