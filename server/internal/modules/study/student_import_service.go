package study

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/table"
)

const emailColumn = "email"

var errImportRejected = errors.New("student import has errors")

type StudentImportService struct {
	tenants database.TenantTx
	gate    CourseGate
	roster  RosterStore
	people  identity.People
	log     audit.Log
}

func NewStudentImportService(tenants database.TenantTx, acc access.Service) StudentImportService {
	return StudentImportService{tenants: tenants, gate: CourseGate{access: acc}}
}

func (s StudentImportService) Import(ctx context.Context, p auth.Principal, courseID string, r io.Reader) (StudentImport, error) {
	t, err := table.Read(r)
	if err != nil {
		return StudentImport{}, fmt.Errorf("%w: %w", ErrBadImport, err)
	}
	column := slices.Index(t.Header, emailColumn)
	if column < 0 {
		return StudentImport{}, fmt.Errorf("%w: column %s is missing", ErrBadImport, emailColumn)
	}
	report := StudentImport{Errors: []ImportError{}}
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, st, err := s.gate.Visible(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		if !st.Manage && st.Role != RoleTeacher {
			return access.ErrForbidden
		}
		if c.ArchivedAt != nil {
			return ErrArchived
		}
		for _, row := range t.Rows {
			code, err := s.add(ctx, tx, p, c, row, column, &report)
			if err != nil {
				return err
			}
			if code != "" {
				report.Errors = append(report.Errors, ImportError{Line: row.Line, Code: code})
			}
		}
		if len(report.Errors) > 0 {
			return errImportRejected
		}
		return nil
	})
	if errors.Is(err, errImportRejected) {
		report.Added, report.Existing = 0, 0
		return report, nil
	}
	return report, err
}

func (s StudentImportService) add(ctx context.Context, tx pgx.Tx, p auth.Principal, c Course, row table.Row, column int, report *StudentImport) (string, error) {
	if column >= len(row.Values) || row.Values[column] == "" {
		return ImportMissingEmail, nil
	}
	ids, err := s.people.IDsByEmail(ctx, tx, row.Values[column])
	if err != nil {
		return "", err
	}
	switch len(ids) {
	case 0:
		return ImportUnknownUser, nil
	case 1:
	default:
		return ImportAmbiguousUser, nil
	}
	added, err := s.roster.Add(ctx, tx, ListStudents, c.ID, ids[0])
	if err != nil {
		return "", err
	}
	if !added {
		report.Existing++
		return "", nil
	}
	report.Added++
	return "", s.log.Record(ctx, tx, p.UserID, audit.Change{
		Action: ActionStudentAdded, TargetKind: "study.course", TargetID: c.ID,
		Details: map[string]any{"title": c.Title, "member": ids[0], "source": "import"},
	})
}
