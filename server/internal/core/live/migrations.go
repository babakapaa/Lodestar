package live

import (
	"embed"
	"io/fs"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrations() database.Migrations {
	sub, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		panic(err)
	}
	return database.Migrations{Owner: "core.live", FS: sub}
}
