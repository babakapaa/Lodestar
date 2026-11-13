package league

import (
	"cmp"
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

type boardRow struct {
	UserID    string
	Division  int
	Points    int
	Growth    int
	TopGrowth bool
	Display   Display
	Nickname  string
}

type board struct {
	rows   []boardRow
	people map[string]identity.Person
}

func (s Service) board(ctx context.Context, tx pgx.Tx, at synced) (board, error) {
	totals, err := s.entries.Totals(ctx, tx, at.season.ID)
	if err != nil {
		return board{}, err
	}
	before := map[string]int{}
	previous, err := s.seasons.Number(ctx, tx, at.league.ID, at.season.Number-1)
	if err != nil {
		return board{}, err
	}
	if previous != nil {
		if before, err = s.entries.Totals(ctx, tx, previous.ID); err != nil {
			return board{}, err
		}
	}
	players, err := s.players.All(ctx, tx, at.league.ID)
	if err != nil {
		return board{}, err
	}
	list, err := s.people.ByIDs(ctx, tx, at.roster.Members)
	if err != nil {
		return board{}, err
	}
	b := board{people: make(map[string]identity.Person, len(list))}
	for _, person := range list {
		b.people[person.UserID] = person
	}
	best := map[int]int{}
	for _, userID := range at.roster.Members {
		player, ok := players[userID]
		if !ok {
			player = Player{UserID: userID, Display: DisplayName}
		}
		row := boardRow{
			UserID: userID, Division: at.league.Divisions.clamp(player.Division), Points: totals[userID],
			Growth: totals[userID] - before[userID], Display: player.Display, Nickname: player.Nickname,
		}
		if row.Growth > best[row.Division] {
			best[row.Division] = row.Growth
		}
		b.rows = append(b.rows, row)
	}
	for i, row := range b.rows {
		b.rows[i].TopGrowth = row.Growth > 0 && row.Growth == best[row.Division]
	}
	return b, nil
}

func (b board) find(userID string) (boardRow, bool) {
	for _, row := range b.rows {
		if row.UserID == userID {
			return row, true
		}
	}
	return boardRow{}, false
}

func (b board) person(userID string) identity.Person {
	if person, ok := b.people[userID]; ok {
		return person
	}
	return identity.Person{UserID: userID}
}

func (b board) table(division int, me string) []TableRow {
	var rows []boardRow
	for _, row := range b.rows {
		if row.Division == division {
			rows = append(rows, row)
		}
	}
	slices.SortFunc(rows, func(a, c boardRow) int {
		return cmp.Or(cmp.Compare(c.Points, a.Points), cmp.Compare(a.UserID, c.UserID))
	})
	out := make([]TableRow, 0, len(rows))
	for i, row := range rows {
		place := i + 1
		if i > 0 && row.Points == rows[i-1].Points {
			place = out[i-1].Place
		}
		out = append(out, TableRow{
			Place: place, Name: b.shown(row), Display: b.display(row), Points: row.Points,
			Growth: row.Growth, Me: row.UserID == me, TopGrowth: row.TopGrowth,
		})
	}
	return out
}

func (b board) display(row boardRow) Display {
	if row.Display == DisplayName && b.people[row.UserID].DisplayName == "" {
		return DisplayHidden
	}
	return row.Display
}

func (b board) shown(row boardRow) string {
	switch b.display(row) {
	case DisplayNickname:
		return row.Nickname
	case DisplayName:
		return b.people[row.UserID].DisplayName
	}
	return ""
}
