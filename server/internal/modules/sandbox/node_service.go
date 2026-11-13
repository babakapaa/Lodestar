package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Nodes(ctx context.Context, p auth.Principal) ([]Node, error) {
	var list []Node
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		var err error
		list, err = s.nodes.All(ctx, tx)
		return err
	})
	return list, err
}

func (s Service) CreateNode(ctx context.Context, p auth.Principal, in NodeInput) (NodeCreated, error) {
	in, err := in.normalized()
	if err != nil {
		return NodeCreated{}, err
	}
	token, err := NewNodeToken(p.TenantID)
	if err != nil {
		return NodeCreated{}, err
	}
	var out NodeCreated
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		n, err := s.nodes.Insert(ctx, tx, in, token.Hash())
		if err != nil {
			return err
		}
		out = NodeCreated{Node: n, Token: token.String()}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionNode, TargetKind: "sandbox.node", TargetID: n.ID, Details: in})
	})
	return out, err
}

func (s Service) DeleteNode(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		if err := s.nodes.Delete(ctx, tx, id); err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionNodeGone, TargetKind: "sandbox.node", TargetID: id})
	})
}
