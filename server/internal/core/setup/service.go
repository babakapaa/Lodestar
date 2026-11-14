package setup

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants  database.TenantTx
	access   access.Service
	settings SettingsWriter
	notifier module.Notifier
	pilots   []module.PilotSource
	sources  []module.ReadinessSource
	probes   []Probe
	hardware []module.HardwareSource
	store    Store
	log      audit.Log
}

func NewService(tenants database.TenantTx, acc access.Service, settings SettingsWriter, notifier module.Notifier, pilots []module.PilotSource, sources ...module.ReadinessSource) Service {
	return Service{tenants: tenants, access: acc, settings: settings, notifier: notifier, pilots: pilots, sources: sources}
}

func (s Service) WithHardware(sources []module.HardwareSource, probes ...Probe) Service {
	s.hardware, s.probes = sources, probes
	return s
}

func (s Service) Hardware(ctx context.Context, p auth.Principal) (HardwareReport, error) {
	var out HardwareReport
	err := s.run(ctx, p, func(tx pgx.Tx, _ State) error {
		if _, err := s.editable(ctx, tx, p, SectionDeploy); err != nil {
			return err
		}
		var err error
		out, err = s.store.Hardware(ctx, tx)
		return err
	})
	return out, err
}

func (s Service) CheckHardware(ctx context.Context, p auth.Principal) (HardwareReport, error) {
	if _, err := s.Hardware(ctx, p); err != nil {
		return HardwareReport{}, err
	}
	checks := hardwareChecker{probes: s.probes}.run(ctx)
	var out HardwareReport
	err := s.run(ctx, p, func(tx pgx.Tx, _ State) error {
		for _, src := range s.hardware {
			found, err := src.Hardware(ctx, tx, p.TenantID)
			if err != nil {
				return err
			}
			checks = append(checks, found...)
		}
		var err error
		out, err = s.store.PutHardware(ctx, tx, checks)
		return err
	})
	return out, err
}

func (s Service) Begin(ctx context.Context, tenantID string) error {
	return s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		return s.store.Begin(ctx, tx)
	})
}

func (s Service) Allowed(ctx context.Context, p auth.Principal) (bool, error) {
	var ok bool
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		st, found, err := s.store.State(ctx, tx)
		if err != nil || !found || (st.LaunchedAt != nil && len(st.Pilot) == 0) {
			ok = err == nil
			return err
		}
		if st.LaunchedAt != nil {
			ok, err = s.inPilot(ctx, tx, p, st)
			return err
		}
		view, err := s.status(ctx, tx, p)
		ok = view.Manages || len(view.Assigned) > 0
		return err
	})
	return ok, err
}

func (s Service) Status(ctx context.Context, p auth.Principal) (StatusView, error) {
	var out StatusView
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		out, err = s.status(ctx, tx, p)
		return err
	})
	return out, err
}

func (s Service) status(ctx context.Context, tx pgx.Tx, p auth.Principal) (StatusView, error) {
	out := StatusView{Assigned: []string{}}
	st, found, err := s.store.State(ctx, tx)
	if err != nil {
		return out, err
	}
	out.Launched = !found || st.LaunchedAt != nil
	out.Included = out.Launched
	if out.Launched && found && len(st.Pilot) > 0 {
		out.Pilot = true
		if out.Manages, err = s.allowed(ctx, tx, p, Manage); err != nil {
			return out, err
		}
		out.Included, err = s.inPilot(ctx, tx, p, st)
		return out, err
	}
	if out.Launched {
		return out, nil
	}
	if out.Manages, err = s.allowed(ctx, tx, p, Manage); err != nil {
		return out, err
	}
	sections, err := s.store.Sections(ctx, tx)
	if err != nil {
		return out, err
	}
	for _, sec := range Sections {
		if a := sections[sec.Key].AssigneeID; a != nil && *a == p.UserID {
			out.Assigned = append(out.Assigned, sec.Key)
		}
	}
	return out, nil
}

func (s Service) Overview(ctx context.Context, p auth.Principal) (Overview, error) {
	var out Overview
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		view, err := s.status(ctx, tx, p)
		if err != nil {
			return err
		}
		legal, err := s.allowed(ctx, tx, p, Legal)
		if err != nil {
			return err
		}
		signer, err := s.allowed(ctx, tx, p, Sign)
		if err != nil {
			return err
		}
		if !view.Manages && !legal && !signer && len(view.Assigned) == 0 {
			return access.ErrForbidden
		}
		out = Overview{
			Launched: st.LaunchedAt != nil, LaunchedAt: st.LaunchedAt, ProtocolVersion: st.ProtocolVersion,
			SignedBy: st.SignedBy, SignedAt: st.SignedAt, Sections: []SectionState{}, Legal: []LegalAnswer{}, Pilot: append([]string{}, st.Pilot...),
		}
		sections, err := s.sectionStates(ctx, tx)
		if err != nil {
			return err
		}
		for _, sec := range sections {
			if view.Manages || slices.Contains(view.Assigned, sec.Key) || (sec.Key == SectionLegal && (legal || signer)) {
				out.Sections = append(out.Sections, sec)
			}
		}
		if view.Manages || legal || signer || slices.Contains(view.Assigned, SectionLegal) {
			if out.Legal, err = s.answers(ctx, tx); err != nil {
				return err
			}
		}
		if view.Manages {
			out.Problems, err = s.problems(ctx, tx, p.TenantID, st)
		}
		return err
	})
	return out, err
}

func (s Service) Assign(ctx context.Context, p auth.Principal, key string, in AssignInput) (SectionState, error) {
	var out SectionState
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		if err := s.require(ctx, tx, p, Manage); err != nil {
			return err
		}
		sec, err := s.section(ctx, tx, key)
		if err != nil {
			return err
		}
		if in.DueOn != nil {
			if _, err := time.Parse(time.DateOnly, *in.DueOn); err != nil {
				return ErrInvalidAssignee
			}
		}
		if in.UserID != nil && key == SectionLegal {
			ok, err := s.access.AllowedInTenant(ctx, tx, auth.Principal{TenantID: p.TenantID, UserID: *in.UserID}, string(Legal))
			if err != nil {
				return err
			}
			if !ok {
				return ErrInvalidAssignee
			}
		}
		before := sec.AssigneeID
		sec.AssigneeID, sec.DueOn = in.UserID, in.DueOn
		if in.UserID == nil {
			sec.DueOn = nil
		}
		if out, err = s.save(ctx, tx, p, sec); err != nil {
			return err
		}
		if in.UserID == nil || *in.UserID == p.UserID || (before != nil && *before == *in.UserID) {
			return nil
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyAssigned, UserIDs: []string{*in.UserID}, Link: sectionLink(key),
		})
	})
	return out, err
}

func (s Service) SetSetting(ctx context.Context, p auth.Principal, key, setting string, value json.RawMessage) error {
	if SectionOf(setting) != key {
		return ErrWrongSection
	}
	return s.run(ctx, p, func(tx pgx.Tx, st State) error {
		sec, err := s.editable(ctx, tx, p, key)
		if err != nil {
			return err
		}
		if value == nil || string(value) == "null" {
			err = s.settings.Clear(ctx, tx, p.UserID, setting)
		} else {
			err = s.settings.Put(ctx, tx, p.UserID, setting, value)
		}
		if err != nil {
			return err
		}
		return s.touch(ctx, tx, p, sec)
	})
}

func (s Service) Submit(ctx context.Context, p auth.Principal, key string) (SectionState, error) {
	var out SectionState
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		sec, err := s.editable(ctx, tx, p, key)
		if err != nil {
			return err
		}
		if sec.Status == InReview {
			return ErrInvalidStatus
		}
		sec.Status = InReview
		out, err = s.save(ctx, tx, p, sec)
		return err
	})
	return out, err
}

func (s Service) Review(ctx context.Context, p auth.Principal, key string, in ReviewInput) (SectionState, error) {
	var out SectionState
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		if err := s.require(ctx, tx, p, Manage); err != nil {
			return err
		}
		sec, err := s.section(ctx, tx, key)
		if err != nil {
			return err
		}
		sec.Note = strings.TrimSpace(in.Note)
		sec.Status = InProgress
		if in.Approve {
			if key == SectionLegal {
				if err := s.legalClosed(ctx, tx); err != nil {
					return err
				}
			}
			sec.Status = Approved
		}
		out, err = s.save(ctx, tx, p, sec)
		return err
	})
	return out, err
}

func (s Service) Answer(ctx context.Context, p auth.Principal, question string, in AnswerInput) (LegalAnswer, error) {
	if !knownQuestion(question) {
		return LegalAnswer{}, ErrUnknownQuestion
	}
	answer, basis := strings.TrimSpace(in.Answer), strings.TrimSpace(in.Basis)
	if answer == "" || basis == "" || len([]rune(answer)) > 4000 || len([]rune(basis)) > 1000 {
		return LegalAnswer{}, ErrInvalidAnswer
	}
	var out LegalAnswer
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		sec, err := s.legalEditor(ctx, tx, p)
		if err != nil {
			return err
		}
		if err := s.store.PutAnswer(ctx, tx, question, answer, basis, p.UserID); err != nil {
			return err
		}
		if err := s.store.Unsign(ctx, tx); err != nil {
			return err
		}
		if err := s.touch(ctx, tx, p, sec); err != nil {
			return err
		}
		if err := s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionAnswered, TargetKind: "setup", TargetID: question}); err != nil {
			return err
		}
		out, err = s.answer(ctx, tx, question)
		return err
	})
	return out, err
}

func (s Service) ApproveAnswer(ctx context.Context, p auth.Principal, question string) (LegalAnswer, error) {
	if !knownQuestion(question) {
		return LegalAnswer{}, ErrUnknownQuestion
	}
	var out LegalAnswer
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		if err := s.require(ctx, tx, p, Legal); err != nil {
			return err
		}
		found, err := s.store.ApproveAnswer(ctx, tx, question, p.UserID)
		if err != nil {
			return err
		}
		if !found {
			return ErrInvalidAnswer
		}
		if err := s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionAnswerApproved, TargetKind: "setup", TargetID: question}); err != nil {
			return err
		}
		out, err = s.answer(ctx, tx, question)
		return err
	})
	return out, err
}

func (s Service) Protocol(ctx context.Context, p auth.Principal) (Protocol, error) {
	var out Protocol
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		ok, err := s.anyOf(ctx, tx, p, Manage, Legal, Sign)
		if err != nil {
			return err
		}
		if !ok {
			return access.ErrForbidden
		}
		answers, err := s.answers(ctx, tx)
		if err != nil {
			return err
		}
		out = Protocol{Version: st.ProtocolVersion, SignedBy: st.SignedBy, SignedAt: st.SignedAt, Answers: answers}
		return nil
	})
	return out, err
}

func (s Service) Sign(ctx context.Context, p auth.Principal) (Overview, error) {
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		if err := s.require(ctx, tx, p, Sign); err != nil {
			return err
		}
		if err := s.legalClosed(ctx, tx); err != nil {
			return err
		}
		version, err := s.store.Sign(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		raw, _ := json.Marshal(version)
		return s.log.Record(ctx, tx, p.UserID, audit.Change{
			Action: ActionSigned, TargetKind: "setup", TargetID: "protocol", Details: map[string]json.RawMessage{"version": raw},
		})
	})
	if err != nil {
		return Overview{}, err
	}
	return s.Overview(ctx, p)
}

func (s Service) Launch(ctx context.Context, p auth.Principal, in LaunchInput) (Overview, error) {
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		if err := s.require(ctx, tx, p, Manage); err != nil {
			return err
		}
		problems, err := s.problems(ctx, tx, p.TenantID, st)
		if err != nil {
			return err
		}
		if len(problems) > 0 {
			return ErrNotReady
		}
		if err := s.validPilot(ctx, tx, p.TenantID, in.Pilot); err != nil {
			return err
		}
		if err := s.store.Launch(ctx, tx, p.UserID, in.Pilot); err != nil {
			return err
		}
		if err := s.store.ClearAssignees(ctx, tx); err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionLaunched, TargetKind: "setup", TargetID: "launch"})
	})
	if err != nil {
		return Overview{}, err
	}
	return s.Overview(ctx, p)
}

func (s Service) run(ctx context.Context, p auth.Principal, fn func(pgx.Tx, State) error) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		st, found, err := s.store.State(ctx, tx)
		if err != nil {
			return err
		}
		if !found {
			return ErrNoSetup
		}
		return fn(tx, st)
	})
}

func (s Service) problems(ctx context.Context, tx pgx.Tx, tenantID string, st State) ([]module.Problem, error) {
	if st.LaunchedAt != nil {
		return []module.Problem{}, nil
	}
	out := []module.Problem{}
	sections, err := s.store.Sections(ctx, tx)
	if err != nil {
		return nil, err
	}
	for _, sec := range Sections {
		if sec.Required && sections[sec.Key].Status != Approved {
			out = append(out, module.Problem{Section: sec.Key, Code: "section_not_approved", Count: 1})
		}
	}
	open, err := s.openAnswers(ctx, tx)
	if err != nil {
		return nil, err
	}
	if open > 0 {
		out = append(out, module.Problem{Section: SectionLegal, Code: "legal_open", Count: open})
	}
	if st.SignedAt == nil {
		out = append(out, module.Problem{Section: SectionLaunch, Code: "protocol_unsigned", Count: 1})
	}
	var units int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM core.units WHERE archived_at IS NULL`).Scan(&units); err != nil {
		return nil, err
	}
	if units == 0 {
		out = append(out, module.Problem{Section: SectionUniversity, Code: "no_units", Count: 1})
	}
	if len(s.probes) > 0 {
		report, err := s.store.Hardware(ctx, tx)
		if err != nil {
			return nil, err
		}
		switch {
		case report.CheckedAt == nil:
			out = append(out, module.Problem{Section: SectionDeploy, Code: "hardware_unchecked", Count: 1})
		case report.failed() > 0:
			out = append(out, module.Problem{Section: SectionDeploy, Code: "hardware_failed", Count: report.failed()})
		}
	}
	for _, src := range s.sources {
		found, err := src.Problems(ctx, tx, tenantID)
		if err != nil {
			return nil, err
		}
		out = append(out, found...)
	}
	return out, nil
}

func (s Service) openAnswers(ctx context.Context, tx pgx.Tx) (int, error) {
	answers, err := s.store.Legal(ctx, tx)
	if err != nil {
		return 0, err
	}
	open := 0
	for _, q := range LegalQuestions {
		if answers[q.Key].ApprovedAt == nil {
			open++
		}
	}
	return open, nil
}

func (s Service) legalClosed(ctx context.Context, tx pgx.Tx) error {
	open, err := s.openAnswers(ctx, tx)
	if err != nil {
		return err
	}
	if open > 0 {
		return ErrLegalOpen
	}
	return nil
}

func (s Service) answers(ctx context.Context, tx pgx.Tx) ([]LegalAnswer, error) {
	stored, err := s.store.Legal(ctx, tx)
	if err != nil {
		return nil, err
	}
	out := make([]LegalAnswer, len(LegalQuestions))
	for i, q := range LegalQuestions {
		a := stored[q.Key]
		a.Question, a.Title = q.Key, q.Title
		out[i] = a
	}
	return out, nil
}

func (s Service) answer(ctx context.Context, tx pgx.Tx, question string) (LegalAnswer, error) {
	all, err := s.answers(ctx, tx)
	if err != nil {
		return LegalAnswer{}, err
	}
	for _, a := range all {
		if a.Question == question {
			return a, nil
		}
	}
	return LegalAnswer{}, ErrUnknownQuestion
}

func (s Service) sectionStates(ctx context.Context, tx pgx.Tx) ([]SectionState, error) {
	stored, err := s.store.Sections(ctx, tx)
	if err != nil {
		return nil, err
	}
	keys := map[string][]string{}
	for _, d := range s.settings.Catalog() {
		sec := SectionOf(d.Key)
		keys[sec] = append(keys[sec], d.Key)
	}
	out := make([]SectionState, len(Sections))
	for i, sec := range Sections {
		st, ok := stored[sec.Key]
		if !ok {
			st = SectionState{Key: sec.Key, Status: NotStarted}
		}
		st.Required = sec.Required
		st.Settings = keys[sec.Key]
		if st.Settings == nil {
			st.Settings = []string{}
		}
		out[i] = st
	}
	return out, nil
}

func (s Service) section(ctx context.Context, tx pgx.Tx, key string) (SectionState, error) {
	all, err := s.sectionStates(ctx, tx)
	if err != nil {
		return SectionState{}, err
	}
	for _, sec := range all {
		if sec.Key == key {
			return sec, nil
		}
	}
	return SectionState{}, ErrUnknownSection
}

func (s Service) editable(ctx context.Context, tx pgx.Tx, p auth.Principal, key string) (SectionState, error) {
	sec, err := s.section(ctx, tx, key)
	if err != nil {
		return SectionState{}, err
	}
	if sec.AssigneeID != nil && *sec.AssigneeID == p.UserID {
		return sec, nil
	}
	return sec, s.require(ctx, tx, p, Manage)
}

func (s Service) legalEditor(ctx context.Context, tx pgx.Tx, p auth.Principal) (SectionState, error) {
	sec, err := s.section(ctx, tx, SectionLegal)
	if err != nil {
		return SectionState{}, err
	}
	if sec.AssigneeID != nil && *sec.AssigneeID == p.UserID {
		return sec, nil
	}
	ok, err := s.anyOf(ctx, tx, p, Manage, Legal)
	if err != nil {
		return SectionState{}, err
	}
	if !ok {
		return SectionState{}, access.ErrForbidden
	}
	return sec, nil
}

func (s Service) touch(ctx context.Context, tx pgx.Tx, p auth.Principal, sec SectionState) error {
	if sec.Status == InProgress {
		return nil
	}
	sec.Status = InProgress
	_, err := s.save(ctx, tx, p, sec)
	return err
}

func (s Service) save(ctx context.Context, tx pgx.Tx, p auth.Principal, sec SectionState) (SectionState, error) {
	if err := s.store.PutSection(ctx, tx, sec, p.UserID); err != nil {
		return SectionState{}, err
	}
	raw, _ := json.Marshal(sec.Status)
	if err := s.log.Record(ctx, tx, p.UserID, audit.Change{
		Action: ActionSectionChanged, TargetKind: "setup", TargetID: sec.Key, Details: map[string]json.RawMessage{"status": raw},
	}); err != nil {
		return SectionState{}, err
	}
	return s.section(ctx, tx, sec.Key)
}

func (s Service) anyOf(ctx context.Context, tx pgx.Tx, p auth.Principal, perms ...access.Permission) (bool, error) {
	for _, perm := range perms {
		ok, err := s.allowed(ctx, tx, p, perm)
		if err != nil || ok {
			return ok, err
		}
	}
	return false, nil
}

func (s Service) require(ctx context.Context, tx pgx.Tx, p auth.Principal, perm access.Permission) error {
	return s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), perm)
}

func (s Service) allowed(ctx context.Context, tx pgx.Tx, p auth.Principal, perm access.Permission) (bool, error) {
	err := s.require(ctx, tx, p, perm)
	if errors.Is(err, access.ErrForbidden) {
		return false, nil
	}
	return err == nil, err
}
