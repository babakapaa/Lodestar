package study

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type FileLink struct {
	gate CourseGate
}

func (l FileLink) Readable(ctx context.Context, tx pgx.Tx, p auth.Principal, fileID string) (bool, error) {
	rows, err := tx.Query(ctx, `
		SELECT a.course_id::text, a.published_at IS NOT NULL, NULL::text FROM study.assignment_files f
		JOIN study.assignments a ON a.id = f.assignment_id WHERE f.file_id::text = $1
		UNION ALL
		SELECT a.course_id::text, true, s.student_id::text FROM study.submission_files f
		JOIN study.submissions s ON s.id = f.submission_id
		JOIN study.assignments a ON a.id = s.assignment_id WHERE f.file_id::text = $1`, fileID)
	if err != nil {
		return false, err
	}
	type use struct {
		course    string
		published bool
		student   *string
	}
	uses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (use, error) {
		var u use
		err := row.Scan(&u.course, &u.published, &u.student)
		return u, err
	})
	if err != nil {
		return false, err
	}
	for _, u := range uses {
		if u.student != nil && *u.student == p.UserID {
			return true, nil
		}
		_, st, err := l.gate.Visible(ctx, tx, p, u.course)
		if errors.Is(err, ErrCourseNotFound) {
			continue
		}
		if err != nil {
			return false, err
		}
		if st.Teaches() || (u.student == nil && u.published) {
			return true, nil
		}
	}
	return false, nil
}

func (FileLink) InUse(ctx context.Context, tx pgx.Tx, fileID string) (bool, error) {
	return exists(ctx, tx, `
		SELECT 1 FROM study.assignment_files WHERE file_id::text = $1
		UNION ALL
		SELECT 1 FROM study.submission_files WHERE file_id::text = $1`, fileID)
}
