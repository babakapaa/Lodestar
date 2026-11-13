package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Estimate(ctx context.Context, p auth.Principal, in EstimateInput) (Estimate, error) {
	in, err := in.normalized()
	if err != nil {
		return Estimate{}, err
	}
	var e Estimate
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, err := s.current(ctx, tx)
		if err != nil {
			return err
		}
		tpl, err := s.templates.ByID(ctx, tx, in.TemplateID)
		if err != nil {
			return err
		}
		if tpl.Archived {
			return ErrTemplateNotFound
		}
		e = Estimate{MaxHours: in.MaxHours, Starter: tpl.Starter, Options: []EstimateOption{}}
		if tpl.Starter {
			e.MaxHours = min(e.MaxHours, starterHours)
		}
		if e.Available, err = s.available(ctx, tx, t, p.UserID, in.TeamID); err != nil {
			return err
		}
		tr, err := s.tariff(ctx, tx)
		if err != nil {
			return err
		}
		if tr.discount > 0 && tr.from != tr.to {
			_, zone := s.location(ctx, tx)
			e.Night = &NightWindow{From: tr.from, To: tr.to, Discount: tr.discount, TimeZone: zone}
		}
		return s.options(ctx, tx, tr, &e)
	})
	return e, err
}

func (s Service) options(ctx context.Context, tx pgx.Tx, tr tariff, e *Estimate) error {
	classes, err := s.classes.All(ctx, tx)
	if err != nil {
		return err
	}
	load, err := s.nodes.Load(ctx, tx)
	if err != nil {
		return err
	}
	var best *EstimateOption
	for _, c := range classes {
		if c.Archived {
			continue
		}
		rate := c.Rate
		if e.Starter {
			rate = 0
		}
		cost := hours(e.MaxHours * rate)
		o := EstimateOption{
			ClassID: c.ID, ClassName: c.Name, MemoryGB: c.MemoryGB, Rate: c.Rate, Cost: cost,
			NightCost: hours(cost * (1 - float64(tr.discount)/100)), FreeGPUs: load[c.ID].free, Queued: load[c.ID].queued,
			Enough: cost <= e.Available,
		}
		if e.Night == nil {
			o.NightCost = cost
		}
		e.Options = append(e.Options, o)
		if o.Enough && (best == nil || o.better(*best)) {
			best = &o
		}
	}
	if best != nil {
		e.Suggested = &best.ClassID
	}
	return nil
}

func (s Service) available(ctx context.Context, tx pgx.Tx, t Term, userID, teamID string) (float64, error) {
	ownerKind, ownerID := OwnerUser, userID
	var lines []QuotaLine
	var err error
	if teamID != "" {
		ref, err := s.team(ctx, tx, userID, teamID)
		if err != nil {
			return 0, err
		}
		ownerKind, ownerID = OwnerTeam, ref.ID
		lines, err = s.ledger.Lines(ctx, tx, t.ID, OwnerTeam, ref.ID)
		if err != nil {
			return 0, err
		}
	} else if lines, err = s.userLines(ctx, tx, t, userID); err != nil {
		return 0, err
	}
	reserved, err := s.runs.Reserved(ctx, tx, t.ID, ownerKind, ownerID)
	if err != nil {
		return 0, err
	}
	return hours(total(lines) - reserved), nil
}
