package module

import "github.com/babakapaa/lodestar/server/internal/platform/database"

type Migrating interface {
	Migrations() database.Migrations
}
