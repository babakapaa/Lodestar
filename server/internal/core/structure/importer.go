package structure

import (
	"context"
	"errors"
	"io"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

var errImportFailed = errors.New("import has errors")

type rowFailed struct {
	code string
}

func (e rowFailed) Error() string {
	return "import row failed: " + e.code
}

func eachRow(ctx context.Context, tx pgx.Tx, rows []importRow, report *ImportReport, fn func(pgx.Tx, importRow) (string, error)) error {
	for _, row := range rows {
		err := pgx.BeginFunc(ctx, tx, func(sp pgx.Tx) error {
			code, err := fn(sp, row)
			if err == nil && code != "" {
				return rowFailed{code: code}
			}
			return err
		})
		var failed rowFailed
		if errors.As(err, &failed) {
			report.Errors = append(report.Errors, ImportError{Line: row.line, Code: failed.code})
			continue
		}
		if err != nil {
			return err
		}
	}
	if len(report.Errors) > 0 {
		return errImportFailed
	}
	return nil
}

func (s Service) ImportUnits(ctx context.Context, p auth.Principal, r io.Reader) (ImportReport, error) {
	rows, err := readTable(r, "code", "parent_code", "kind", "name")
	if err != nil {
		return ImportReport{}, err
	}
	var report ImportReport
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), ManageStructure); err != nil {
			return err
		}
		return eachRow(ctx, tx, rows, &report, func(sp pgx.Tx, row importRow) (string, error) {
			return s.importUnit(ctx, sp, p, row, &report), nil
		})
	})
	return finishImport(report, err)
}

func (s Service) ImportMembers(ctx context.Context, p auth.Principal, r io.Reader) (ImportReport, error) {
	rows, err := readTable(r, "unit_code", "email")
	if err != nil {
		return ImportReport{}, err
	}
	var report ImportReport
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), ManageMembers); err != nil {
			return err
		}
		return eachRow(ctx, tx, rows, &report, func(sp pgx.Tx, row importRow) (string, error) {
			return s.importMember(ctx, sp, p, row, &report)
		})
	})
	return finishImport(report, err)
}

func (s Service) importUnit(ctx context.Context, tx pgx.Tx, p auth.Principal, row importRow, report *ImportReport) string {
	code := row.values["code"]
	in := UnitInput{Kind: row.values["kind"], Name: row.values["name"], Code: &code}
	if code == "" {
		return ImportInvalidUnit
	}
	if parentCode := row.values["parent_code"]; parentCode != "" {
		parent, err := s.units.ByCode(ctx, tx, parentCode)
		if err != nil {
			return ImportMissingParent
		}
		in.ParentID = &parent.ID
	}
	existing, err := s.units.ByCode(ctx, tx, code)
	if errors.Is(err, ErrUnitNotFound) {
		u, err := s.units.Insert(ctx, tx, in)
		if err != nil {
			return ImportInvalidUnit
		}
		report.Created++
		if s.events.Record(ctx, tx, u.ID, &p.UserID, ActionCreated, in) != nil {
			return ImportInvalidUnit
		}
		return ""
	}
	if err != nil {
		return ImportInvalidUnit
	}
	if err := s.checkNoCycle(ctx, tx, existing.ID, in.ParentID); err != nil {
		return ImportCycle
	}
	if _, err := s.units.Update(ctx, tx, existing.ID, in); err != nil {
		return ImportInvalidUnit
	}
	if err := s.units.SetParent(ctx, tx, existing.ID, in.ParentID); err != nil {
		return ImportInvalidUnit
	}
	report.Updated++
	if s.events.Record(ctx, tx, existing.ID, &p.UserID, ActionUpdated, map[string]any{"before": existing, "after": in}) != nil {
		return ImportInvalidUnit
	}
	return ""
}

func (s Service) importMember(ctx context.Context, tx pgx.Tx, p auth.Principal, row importRow, report *ImportReport) (string, error) {
	unit, err := s.units.ByCode(ctx, tx, row.values["unit_code"])
	if errors.Is(err, ErrUnitNotFound) {
		return ImportMissingUnit, nil
	}
	if err != nil {
		return "", err
	}
	ids, err := s.people.IDsByEmail(ctx, tx, row.values["email"])
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
	added, err := s.members.Add(ctx, tx, unit.ID, ids[0])
	if err != nil {
		return "", err
	}
	if added {
		report.Created++
		return "", s.events.Record(ctx, tx, unit.ID, &p.UserID, ActionMemberAdded, map[string]string{"userId": ids[0]})
	}
	return "", nil
}

func finishImport(report ImportReport, err error) (ImportReport, error) {
	if errors.Is(err, errImportFailed) {
		report.Created, report.Updated = 0, 0
		return report, nil
	}
	if err != nil {
		return ImportReport{}, err
	}
	if report.Errors == nil {
		report.Errors = []ImportError{}
	}
	return report, nil
}
