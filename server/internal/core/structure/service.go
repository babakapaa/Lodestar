package structure

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const historyPage = 50

type Service struct {
	tenants database.TenantTx
	access  access.Service
	units   UnitStore
	members MemberStore
	events  EventStore
}

func NewService(tenants database.TenantTx, acc access.Service) Service {
	return Service{tenants: tenants, access: acc}
}

func (s Service) Children(ctx context.Context, p auth.Principal, parentID *string, archived bool) ([]Unit, error) {
	var list []Unit
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.units.Children(ctx, tx, parentID, archived)
		return err
	})
	return list, err
}

func (s Service) Get(ctx context.Context, p auth.Principal, id string) (UnitDetails, error) {
	var d UnitDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		path, err := s.units.Path(ctx, tx, id)
		if err != nil {
			return err
		}
		if len(path) == 0 {
			return ErrUnitNotFound
		}
		d = UnitDetails{Unit: path[len(path)-1], Path: path[:len(path)-1]}
		return nil
	})
	return d, err
}

func (s Service) Create(ctx context.Context, p auth.Principal, in UnitInput) (Unit, error) {
	var u Unit
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.requireParent(ctx, tx, p, in.ParentID, ManageStructure); err != nil {
			return err
		}
		var err error
		if u, err = s.units.Insert(ctx, tx, in); err != nil {
			return err
		}
		return s.events.Record(ctx, tx, u.ID, &p.UserID, ActionCreated, in)
	})
	return u, err
}

func (s Service) Update(ctx context.Context, p auth.Principal, id string, in UnitInput) (Unit, error) {
	var u Unit
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		before, err := s.units.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if err := s.access.Require(ctx, tx, p, UnitScope(id), ManageStructure); err != nil {
			return err
		}
		if u, err = s.units.Update(ctx, tx, id, in); err != nil {
			return err
		}
		return s.events.Record(ctx, tx, id, &p.UserID, ActionUpdated, map[string]any{"before": before, "after": u})
	})
	return u, err
}

func (s Service) Move(ctx context.Context, p auth.Principal, id string, parentID *string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		before, err := s.units.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if err := s.access.Require(ctx, tx, p, UnitScope(id), ManageStructure); err != nil {
			return err
		}
		if err := s.requireParent(ctx, tx, p, parentID, ManageStructure); err != nil {
			return err
		}
		if err := s.checkNoCycle(ctx, tx, id, parentID); err != nil {
			return err
		}
		if err := s.units.SetParent(ctx, tx, id, parentID); err != nil {
			return err
		}
		return s.events.Record(ctx, tx, id, &p.UserID, ActionMoved, map[string]any{"from": before.ParentID, "to": parentID})
	})
}

func (s Service) SetArchived(ctx context.Context, p auth.Principal, id string, archived bool) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		u, err := s.units.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if err := s.access.Require(ctx, tx, p, UnitScope(id), ManageStructure); err != nil {
			return err
		}
		if !archived && u.ParentID != nil {
			parent, err := s.units.ByID(ctx, tx, *u.ParentID)
			if err != nil {
				return err
			}
			if parent.ArchivedAt != nil {
				return ErrArchived
			}
		}
		if err := s.units.SetArchived(ctx, tx, id, archived); err != nil {
			return err
		}
		action := ActionRestored
		if archived {
			action = ActionArchived
		}
		return s.events.Record(ctx, tx, id, &p.UserID, action, struct{}{})
	})
}

func (s Service) Members(ctx context.Context, p auth.Principal, unitID string) ([]Member, error) {
	var list []Member
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.units.ByID(ctx, tx, unitID); err != nil {
			return err
		}
		var err error
		list, err = s.members.List(ctx, tx, unitID)
		return err
	})
	return list, err
}

func (s Service) AddMember(ctx context.Context, p auth.Principal, unitID, userID string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.memberChange(ctx, tx, p, unitID); err != nil {
			return err
		}
		added, err := s.members.Add(ctx, tx, unitID, userID)
		if err != nil || !added {
			return err
		}
		return s.events.Record(ctx, tx, unitID, &p.UserID, ActionMemberAdded, map[string]string{"userId": userID})
	})
}

func (s Service) RemoveMember(ctx context.Context, p auth.Principal, unitID, userID string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.memberChange(ctx, tx, p, unitID); err != nil {
			return err
		}
		removed, err := s.members.Remove(ctx, tx, unitID, userID)
		if err != nil || !removed {
			return err
		}
		return s.events.Record(ctx, tx, unitID, &p.UserID, ActionMemberRemoved, map[string]string{"userId": userID})
	})
}

func (s Service) UserUnits(ctx context.Context, p auth.Principal, userID string) ([]Unit, error) {
	var list []Unit
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.members.UnitsOf(ctx, tx, userID)
		return err
	})
	return list, err
}

func (s Service) History(ctx context.Context, p auth.Principal, unitID string, before int64) ([]Event, error) {
	var list []Event
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.units.ByID(ctx, tx, unitID); err != nil {
			return err
		}
		if err := s.access.Require(ctx, tx, p, UnitScope(unitID), ManageStructure); err != nil {
			return err
		}
		var err error
		list, err = s.events.List(ctx, tx, unitID, before, historyPage)
		return err
	})
	return list, err
}

func (s Service) requireParent(ctx context.Context, tx pgx.Tx, p auth.Principal, parentID *string, perm access.Permission) error {
	if parentID == nil {
		return s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), perm)
	}
	parent, err := s.units.ByID(ctx, tx, *parentID)
	if err != nil {
		return err
	}
	if parent.ArchivedAt != nil {
		return ErrArchived
	}
	return s.access.Require(ctx, tx, p, UnitScope(parent.ID), perm)
}

func (s Service) memberChange(ctx context.Context, tx pgx.Tx, p auth.Principal, unitID string) error {
	u, err := s.units.ByID(ctx, tx, unitID)
	if err != nil {
		return err
	}
	if u.ArchivedAt != nil {
		return ErrArchived
	}
	return s.access.Require(ctx, tx, p, UnitScope(unitID), ManageMembers)
}

func (s Service) checkNoCycle(ctx context.Context, tx pgx.Tx, id string, parentID *string) error {
	if parentID == nil {
		return nil
	}
	inside, err := s.units.IsInside(ctx, tx, *parentID, id)
	if err != nil {
		return err
	}
	if inside {
		return ErrCycle
	}
	return nil
}
