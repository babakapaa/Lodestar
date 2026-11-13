package league

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type PlayerRow struct {
	Person    identity.Person `json:"person"`
	Division  int             `json:"division"`
	Points    int             `json:"points"`
	Growth    int             `json:"growth"`
	TopGrowth bool            `json:"topGrowth"`
	Display   Display         `json:"display"`
	Nickname  string          `json:"nickname"`
}
