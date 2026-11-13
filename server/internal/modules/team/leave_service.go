package team

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) Leave(ctx context.Context, p auth.Principal, id string, in LeaveInput) (LeaveResult, error) {
	in, err := in.normalized()
	if err != nil {
		return LeaveResult{}, err
	}
	var res LeaveResult
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, ts, _, err := s.team(ctx, tx, p, id, true)
		if err != nil {
			return err
		}
		if !t.has(p.UserID) {
			return ErrTeamNotFound
		}
		if freeLeave(ts, time.Now()) {
			res.Left = true
			return s.drop(ctx, tx, t, p.UserID, EventLeft, p.UserID)
		}
		requestID, err := s.leaves.Insert(ctx, tx, t.ID, p.UserID, in.Reason)
		if err != nil {
			return err
		}
		request, err := s.leaves.ByID(ctx, tx, requestID)
		if err != nil {
			return err
		}
		res.Request = &request
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyLeaveRequested, UserIDs: ts.Leads,
			Params: map[string]any{"scope": ts.Title, "team": t.Name, "student": request.Student.DisplayName},
			Link:   teamLink(t.ID),
		})
	})
	return res, err
}

func (s Service) Decide(ctx context.Context, p auth.Principal, requestID string, in DecisionInput) (TeamDetails, error) {
	var d TeamDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		request, err := s.leaves.ByID(ctx, tx, requestID)
		if err != nil {
			return err
		}
		t, ts, role, err := s.team(ctx, tx, p, request.TeamID, true)
		if err != nil || role != RoleLead {
			return ErrRequestNotFound
		}
		status, kind := LeaveRejected, NotifyLeaveRejected
		if in.Approve {
			status, kind = LeaveApproved, NotifyLeaveApproved
		}
		decided, err := s.leaves.Decide(ctx, tx, requestID, status, p.UserID)
		if err != nil {
			return err
		}
		if !decided {
			return ErrRequestNotFound
		}
		if in.Approve && t.has(request.Student.UserID) {
			if err := s.drop(ctx, tx, t, request.Student.UserID, EventRemoved, p.UserID); err != nil {
				return err
			}
		}
		if err := s.log.Record(ctx, tx, p.UserID, audit.Change{
			Action: ActionLeave, TargetKind: "team.team", TargetID: t.ID,
			Details: map[string]any{"student": request.Student.UserID, "approved": in.Approve},
		}); err != nil {
			return err
		}
		if err := s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: kind, UserIDs: []string{request.Student.UserID},
			Params: map[string]any{"scope": ts.Title, "team": t.Name},
			Link:   teamLink(t.ID),
		}); err != nil {
			return err
		}
		if t, err = s.teams.ByID(ctx, tx, t.ID, false); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, t, ts, role)
		return err
	})
	return d, err
}
