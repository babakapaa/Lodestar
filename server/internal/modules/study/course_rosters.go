package study

import (
	"context"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

const RosterCourse = "study.course"

type CourseRosters struct {
	directory structure.Directory
}

func (CourseRosters) Kind() string {
	return RosterCourse
}

func (r CourseRosters) All(ctx context.Context, tx pgx.Tx) ([]module.Roster, error) {
	return r.load(ctx, tx, nil)
}

func (r CourseRosters) One(ctx context.Context, tx pgx.Tx, id string) (module.Roster, error) {
	list, err := r.load(ctx, tx, []string{id})
	if err != nil {
		return module.Roster{}, err
	}
	if len(list) == 0 {
		return module.Roster{}, module.ErrRosterNotFound
	}
	return list[0], nil
}

func (r CourseRosters) load(ctx context.Context, tx pgx.Tx, ids []string) ([]module.Roster, error) {
	rows, err := tx.Query(ctx, `
		SELECT c.id::text, c.title, c.archived_at IS NOT NULL,
		       coalesce((SELECT array_agg(t.user_id::text) FROM study.course_teachers t WHERE t.course_id = c.id), '{}'),
		       coalesce((SELECT array_agg(s.user_id::text) FROM study.course_students s WHERE s.course_id = c.id), '{}'),
		       coalesce((SELECT array_agg(u.unit_id::text) FROM study.course_units u WHERE u.course_id = c.id), '{}')
		FROM study.courses c
		WHERE $1::text[] IS NULL OR c.id::text = ANY($1)
		ORDER BY c.id`, ids)
	if err != nil {
		return nil, fmt.Errorf("course rosters: %w", err)
	}
	list, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (courseRosterRow, error) {
		c := courseRosterRow{roster: module.Roster{Kind: RosterCourse}}
		err := row.Scan(&c.roster.ID, &c.roster.Title, &c.roster.Archived, &c.roster.Leads, &c.students, &c.units)
		return c, err
	})
	if err != nil {
		return nil, fmt.Errorf("course rosters: %w", err)
	}
	var units []string
	for _, c := range list {
		units = append(units, c.units...)
	}
	under, err := r.directory.UsersUnderEach(ctx, tx, units)
	if err != nil {
		return nil, err
	}
	out := make([]module.Roster, 0, len(list))
	for _, c := range list {
		members := slices.Clone(c.students)
		for _, unit := range c.units {
			members = append(members, under[unit]...)
		}
		slices.Sort(members)
		members = slices.Compact(members)
		c.roster.Members = slices.DeleteFunc(members, func(id string) bool { return slices.Contains(c.roster.Leads, id) })
		slices.Sort(c.roster.Leads)
		out = append(out, c.roster)
	}
	return out, nil
}
