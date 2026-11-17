package study

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type DefenseService struct {
	tenants  database.TenantTx
	work     AssignmentService
	settings module.Settings
	notifier module.Notifier
	slots    SlotStore
	defenses DefenseStore
	grades   GradeStore
	drafter  DefenseDrafter
	people   identity.People
	log      audit.Log
	now      func() time.Time
}

func NewDefenseService(k module.Kernel) DefenseService {
	return DefenseService{
		tenants: k.Tenants, work: NewAssignmentService(k.Tenants, k.Access, k.Notifier),
		settings: k.Settings, notifier: k.Notifier, drafter: DefenseDrafter{model: k.Model}, now: time.Now,
	}
}

func (s DefenseService) Slots(ctx context.Context, p auth.Principal, id string) ([]DefenseSlot, error) {
	var list []DefenseSlot
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, _, st, err := s.work.visible(ctx, tx, p, id, false)
		if err != nil {
			return err
		}
		if a.Defense == nil {
			return ErrNoDefense
		}
		if list, err = s.slots.ByAssignment(ctx, tx, id); err != nil {
			return err
		}
		if st.Teaches() {
			return s.fillStudents(ctx, tx, id, list)
		}
		mine, err := s.defenses.Get(ctx, tx, id, p.UserID)
		if err != nil {
			return err
		}
		for i := range list {
			list[i].Mine = mine.SlotID != nil && *mine.SlotID == list[i].ID
		}
		return nil
	})
	return list, err
}

func (s DefenseService) AddSlots(ctx context.Context, p auth.Principal, id string, in SlotsInput) ([]DefenseSlot, error) {
	in, err := in.normalized()
	if err != nil {
		return nil, err
	}
	var list []DefenseSlot
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, err := s.taught(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if list, err = s.slots.Insert(ctx, tx, id, p.UserID, in); err != nil {
			return err
		}
		return s.record(ctx, tx, p, ActionDefenseSlotsAdded, a, map[string]any{"count": len(list)})
	})
	return list, err
}

func (s DefenseService) DeleteSlot(ctx context.Context, p auth.Principal, slotID string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		slot, err := s.slots.ByID(ctx, tx, slotID, true)
		if err != nil {
			return err
		}
		a, c, st, err := s.work.visible(ctx, tx, p, slot.AssignmentID, true)
		if err != nil {
			return ErrSlotNotFound
		}
		if !st.Teaches() {
			return access.ErrForbidden
		}
		if c.ArchivedAt != nil {
			return ErrArchived
		}
		students, err := s.defenses.InSlot(ctx, tx, slotID)
		if err != nil {
			return err
		}
		if err := s.slots.Delete(ctx, tx, slotID); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionDefenseSlotDeleted, a, map[string]any{"startsAt": slot.StartsAt}); err != nil {
			return err
		}
		if len(students) == 0 {
			return nil
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyDefenseCancelled, UserIDs: students,
			Params: map[string]any{"assignment": a.Title, "course": c.Title},
			Link:   assignmentLink(a.ID),
		})
	})
}

func (s DefenseService) Mine(ctx context.Context, p auth.Principal, id string) (MyDefense, error) {
	var my MyDefense
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, c, _, err := s.work.visible(ctx, tx, p, id, false)
		if err != nil {
			return err
		}
		if a.Defense == nil {
			return ErrNoDefense
		}
		if err := s.work.student(ctx, tx, c, p.UserID); err != nil {
			return err
		}
		my, err = s.mine(ctx, tx, a, p.UserID)
		return err
	})
	return my, err
}

func (s DefenseService) Book(ctx context.Context, p auth.Principal, id string, in BookingInput) (MyDefense, error) {
	var my MyDefense
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, d, err := s.own(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if a.Defense != nil && *a.Defense == DefenseSample && !d.Selected {
			return ErrNotChosen
		}
		slot, err := s.slots.ByID(ctx, tx, in.SlotID, true)
		if err != nil {
			return err
		}
		if slot.AssignmentID != a.ID {
			return ErrSlotNotFound
		}
		if !slot.StartsAt.After(s.now()) {
			return ErrSlotStarted
		}
		if d.SlotID != nil && *d.SlotID == slot.ID {
			my, err = s.mine(ctx, tx, a, p.UserID)
			return err
		}
		if !slot.Free() {
			return ErrSlotFull
		}
		if err := s.leave(ctx, tx, d); err != nil {
			return err
		}
		if _, err := s.defenses.Book(ctx, tx, a.ID, p.UserID, &slot.ID); err != nil {
			return err
		}
		my, err = s.mine(ctx, tx, a, p.UserID)
		return err
	})
	return my, err
}

func (s DefenseService) Unbook(ctx context.Context, p auth.Principal, id string) (MyDefense, error) {
	var my MyDefense
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, d, err := s.own(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if err := s.leave(ctx, tx, d); err != nil {
			return err
		}
		if _, err := s.defenses.Book(ctx, tx, a.ID, p.UserID, nil); err != nil {
			return err
		}
		my, err = s.mine(ctx, tx, a, p.UserID)
		return err
	})
	return my, err
}

func (s DefenseService) List(ctx context.Context, p auth.Principal, id string) ([]DefenseRow, error) {
	var rows []DefenseRow
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, c, st, err := s.work.visible(ctx, tx, p, id, false)
		if err != nil {
			return err
		}
		if !st.Teaches() {
			return access.ErrForbidden
		}
		if a.Defense == nil {
			return ErrNoDefense
		}
		studentIDs, err := s.work.gate.StudentIDs(ctx, tx, c.ID)
		if err != nil {
			return err
		}
		people, err := s.people.ByIDs(ctx, tx, studentIDs)
		if err != nil {
			return err
		}
		slots, err := s.slots.ByAssignment(ctx, tx, id)
		if err != nil {
			return err
		}
		bySlot := map[string]DefenseSlot{}
		for _, slot := range slots {
			bySlot[slot.ID] = slot
		}
		defenses, err := s.defenses.ByAssignment(ctx, tx, id)
		if err != nil {
			return err
		}
		byStudent := map[string]Defense{}
		for _, d := range defenses {
			byStudent[d.StudentID] = d
		}
		grades, err := s.grades.ByAssignment(ctx, tx, []string{id}, nil)
		if err != nil {
			return err
		}
		graded := map[string]Grade{}
		for _, g := range grades {
			graded[g.StudentID] = g
		}
		rows = make([]DefenseRow, 0, len(people))
		for _, person := range people {
			d, ok := byStudent[person.UserID]
			row := DefenseRow{Student: person, Status: DefensePlanned}
			if ok {
				row.Selected, row.Status = d.Selected, d.Status
				if d.SlotID != nil {
					slot := bySlot[*d.SlotID]
					row.Slot = &slot
				}
			}
			if g, ok := graded[person.UserID]; ok {
				row.Grade = &g
			}
			rows = append(rows, row)
		}
		return nil
	})
	return rows, err
}

func (s DefenseService) Select(ctx context.Context, p auth.Principal, id, studentID string, in SelectionInput) (Defense, error) {
	var d Defense
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.taughtStudent(ctx, tx, p, id, studentID); err != nil {
			return err
		}
		var err error
		d, err = s.defenses.Select(ctx, tx, id, studentID, in.Selected)
		return err
	})
	return d, err
}

func (s DefenseService) Sample(ctx context.Context, p auth.Principal, id string, in SampleInput) ([]string, error) {
	if in.Count < 1 {
		return nil, ErrInvalidDefense
	}
	var chosen []string
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, err := s.taught(ctx, tx, p, id)
		if err != nil {
			return err
		}
		studentIDs, err := s.work.gate.StudentIDs(ctx, tx, a.CourseID)
		if err != nil {
			return err
		}
		defenses, err := s.defenses.ByAssignment(ctx, tx, id)
		if err != nil {
			return err
		}
		already := map[string]bool{}
		for _, d := range defenses {
			already[d.StudentID] = d.Selected
		}
		var pool []string
		for _, sid := range studentIDs {
			if !already[sid] {
				pool = append(pool, sid)
			}
		}
		rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		chosen = pool[:min(in.Count, len(pool))]
		for _, sid := range chosen {
			if _, err := s.defenses.Select(ctx, tx, id, sid, true); err != nil {
				return err
			}
		}
		return nil
	})
	return chosen, err
}

func (s DefenseService) Get(ctx context.Context, p auth.Principal, id, studentID string) (Defense, error) {
	var d Defense
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.taughtStudent(ctx, tx, p, id, studentID); err != nil {
			return err
		}
		var err error
		d, err = s.defenses.Get(ctx, tx, id, studentID)
		return err
	})
	return d, err
}

func (s DefenseService) Save(ctx context.Context, p auth.Principal, id, studentID string, in DefenseInput) (Defense, error) {
	in, err := in.normalized()
	if err != nil {
		return Defense{}, err
	}
	var d Defense
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.taughtStudent(ctx, tx, p, id, studentID); err != nil {
			return err
		}
		d, err = s.defenses.Save(ctx, tx, id, studentID, in)
		return err
	})
	return d, err
}

func (s DefenseService) Finish(ctx context.Context, p auth.Principal, id, studentID string, in FinishInput) (Defense, error) {
	if !in.valid() {
		return Defense{}, ErrInvalidDefense
	}
	var d Defense
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, err := s.taughtStudent(ctx, tx, p, id, studentID)
		if err != nil {
			return err
		}
		if !a.Published() {
			return ErrNotPublished
		}
		c, err := s.work.gate.courses.ByID(ctx, tx, a.CourseID)
		if err != nil {
			return err
		}
		if d, err = s.defenses.Finish(ctx, tx, id, studentID, p.UserID, in.Status); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionDefenseFinished, a, map[string]any{"student": studentID, "status": in.Status}); err != nil {
			return err
		}
		if in.Grade != nil {
			if _, err := s.work.grade(ctx, tx, p, a, c, studentID, *in.Grade); err != nil {
				return err
			}
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyDefenseHeld, UserIDs: []string{studentID},
			Params: map[string]any{"assignment": a.Title, "course": c.Title},
			Link:   assignmentLink(a.ID),
		})
	})
	return d, err
}

func (s DefenseService) DraftProtocol(ctx context.Context, p auth.Principal, id, studentID string, in DefenseInput) (string, error) {
	if !s.drafter.ready() {
		return "", ErrAIUnavailable
	}
	in, err := in.normalized()
	if err != nil {
		return "", err
	}
	var title string
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, err := s.taughtStudent(ctx, tx, p, id, studentID)
		title = a.Title
		return err
	})
	if err != nil {
		return "", err
	}
	return s.drafter.Protocol(ctx, title, in.Questions, in.Summary)
}

func (s DefenseService) mine(ctx context.Context, tx pgx.Tx, a Assignment, studentID string) (MyDefense, error) {
	d, err := s.defenses.Get(ctx, tx, a.ID, studentID)
	if err != nil {
		return MyDefense{}, err
	}
	recorded, err := s.recorded(ctx, tx)
	if err != nil {
		return MyDefense{}, err
	}
	my := MyDefense{
		Format: *a.Defense, Required: *a.Defense != DefenseSample || d.Selected,
		Recorded: recorded, Status: d.Status,
	}
	if d.SlotID != nil {
		slot, err := s.slots.ByID(ctx, tx, *d.SlotID, false)
		if err != nil {
			return MyDefense{}, err
		}
		slot.Mine = true
		my.Slot = &slot
	}
	if d.Status.finished() {
		my.Protocol = &d
	}
	return my, nil
}

func (s DefenseService) own(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Assignment, Defense, error) {
	a, c, _, err := s.work.visible(ctx, tx, p, id, false)
	if err != nil {
		return Assignment{}, Defense{}, err
	}
	if a.Defense == nil {
		return Assignment{}, Defense{}, ErrNoDefense
	}
	if c.ArchivedAt != nil {
		return Assignment{}, Defense{}, ErrArchived
	}
	if err := s.work.student(ctx, tx, c, p.UserID); err != nil {
		return Assignment{}, Defense{}, err
	}
	d, err := s.defenses.Get(ctx, tx, id, p.UserID)
	if err != nil {
		return Assignment{}, Defense{}, err
	}
	if d.Status.finished() {
		return Assignment{}, Defense{}, ErrDefenseFinished
	}
	return a, d, nil
}

func (s DefenseService) leave(ctx context.Context, tx pgx.Tx, d Defense) error {
	if d.SlotID == nil {
		return nil
	}
	current, err := s.slots.ByID(ctx, tx, *d.SlotID, true)
	if err != nil {
		return err
	}
	if !current.StartsAt.After(s.now()) {
		return ErrSlotStarted
	}
	return nil
}

func (s DefenseService) taught(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Assignment, error) {
	a, c, err := s.work.taught(ctx, tx, p, id)
	if err != nil {
		return Assignment{}, err
	}
	if a.Defense == nil {
		return Assignment{}, ErrNoDefense
	}
	if c.ArchivedAt != nil {
		return Assignment{}, ErrArchived
	}
	return a, nil
}

func (s DefenseService) taughtStudent(ctx context.Context, tx pgx.Tx, p auth.Principal, id, studentID string) (Assignment, error) {
	a, err := s.taught(ctx, tx, p, id)
	if err != nil {
		return Assignment{}, err
	}
	c, err := s.work.gate.courses.ByID(ctx, tx, a.CourseID)
	if err != nil {
		return Assignment{}, err
	}
	return a, s.work.student(ctx, tx, c, studentID)
}

func (s DefenseService) fillStudents(ctx context.Context, tx pgx.Tx, id string, list []DefenseSlot) error {
	defenses, err := s.defenses.ByAssignment(ctx, tx, id)
	if err != nil {
		return err
	}
	bySlot := map[string][]string{}
	var all []string
	for _, d := range defenses {
		if d.SlotID != nil {
			bySlot[*d.SlotID] = append(bySlot[*d.SlotID], d.StudentID)
			all = append(all, d.StudentID)
		}
	}
	people, err := s.people.ByIDs(ctx, tx, all)
	if err != nil {
		return err
	}
	byID := map[string]identity.Person{}
	for _, person := range people {
		byID[person.UserID] = person
	}
	for i := range list {
		list[i].Students = []identity.Person{}
		for _, sid := range bySlot[list[i].ID] {
			list[i].Students = append(list[i].Students, byID[sid])
		}
	}
	return nil
}

func (s DefenseService) recorded(ctx context.Context, tx pgx.Tx) (bool, error) {
	raw, err := s.settings.Value(ctx, tx, DefenseRecording)
	if err != nil {
		return false, err
	}
	var on bool
	if err := json.Unmarshal(raw, &on); err != nil {
		return false, fmt.Errorf("read %s: %w", DefenseRecording, err)
	}
	return on, nil
}

func (s DefenseService) record(ctx context.Context, tx pgx.Tx, p auth.Principal, action string, a Assignment, details map[string]any) error {
	details["title"] = a.Title
	return s.log.Record(ctx, tx, p.UserID, audit.Change{
		Action: action, TargetKind: "study.assignment", TargetID: a.ID, Details: details,
	})
}
