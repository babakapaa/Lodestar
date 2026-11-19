package study

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type AssignmentService struct {
	tenants     database.TenantTx
	gate        CourseGate
	notifier    module.Notifier
	assignments AssignmentStore
	submissions SubmissionStore
	maps        SearchMapStore
	grades      GradeStore
	attachments Attachments
	people      identity.People
	log         audit.Log
	teams       module.TeamDirectory
	help        module.HelpSource
}

func NewAssignmentService(tenants database.TenantTx, acc access.Service, notifier module.Notifier) AssignmentService {
	return AssignmentService{tenants: tenants, gate: CourseGate{access: acc}, notifier: notifier}
}

func (s AssignmentService) WithHelp(help module.HelpSource) AssignmentService {
	s.help = help
	return s
}

func (s AssignmentService) WithTeams(teams module.TeamDirectory) AssignmentService {
	s.teams = teams
	return s
}

func (s AssignmentService) List(ctx context.Context, p auth.Principal, courseID string) ([]AssignmentItem, error) {
	var items []AssignmentItem
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		_, st, err := s.gate.Visible(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		list, err := s.assignments.ByCourse(ctx, tx, courseID, st.Teaches())
		if err != nil {
			return err
		}
		ids := make([]string, len(list))
		for i, a := range list {
			ids[i] = a.ID
		}
		items = make([]AssignmentItem, len(list))
		if st.Teaches() {
			return s.fillCounts(ctx, tx, courseID, list, items)
		}
		return s.fillProgress(ctx, tx, p.UserID, ids, list, items)
	})
	return items, err
}

func (s AssignmentService) fillCounts(ctx context.Context, tx pgx.Tx, courseID string, list []Assignment, items []AssignmentItem) error {
	students, err := s.gate.StudentIDs(ctx, tx, courseID)
	if err != nil {
		return err
	}
	ids := make([]string, len(list))
	for i, a := range list {
		ids[i] = a.ID
	}
	submitted, err := s.submissions.Submitters(ctx, tx, ids)
	if err != nil {
		return err
	}
	grades, err := s.grades.ByAssignment(ctx, tx, ids, nil)
	if err != nil {
		return err
	}
	graded := map[string]int{}
	for _, g := range grades {
		graded[g.AssignmentID]++
	}
	for i, a := range list {
		items[i] = AssignmentItem{Assignment: a, Counts: &Counts{Students: len(students), Submitted: submitted[a.ID], Graded: graded[a.ID]}}
	}
	return nil
}

func (s AssignmentService) fillProgress(ctx context.Context, tx pgx.Tx, userID string, ids []string, list []Assignment, items []AssignmentItem) error {
	last, err := s.submissions.LastSubmitted(ctx, tx, ids, userID)
	if err != nil {
		return err
	}
	grades, err := s.grades.ByAssignment(ctx, tx, ids, &userID)
	if err != nil {
		return err
	}
	byAssignment := map[string]Grade{}
	for _, g := range grades {
		byAssignment[g.AssignmentID] = g
	}
	for i, a := range list {
		progress := &Progress{}
		if at, ok := last[a.ID]; ok {
			progress.SubmittedAt = &at
		}
		if g, ok := byAssignment[a.ID]; ok {
			progress.Grade = &g
		}
		items[i] = AssignmentItem{Assignment: a, Mine: progress}
	}
	return nil
}

func (s AssignmentService) Create(ctx context.Context, p auth.Principal, courseID string, in AssignmentInput) (Assignment, error) {
	in, err := in.normalized()
	if err != nil {
		return Assignment{}, err
	}
	var a Assignment
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, st, err := s.gate.Visible(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		if !st.Teaches() {
			return access.ErrForbidden
		}
		if c.ArchivedAt != nil {
			return ErrArchived
		}
		if err := s.attachments.Check(ctx, tx, p.UserID, in.Files); err != nil {
			return err
		}
		if a, err = s.assignments.Insert(ctx, tx, courseID, p.UserID, in); err != nil {
			return err
		}
		if err := s.assignments.SetFiles(ctx, tx, a.ID, in.Files); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionAssignmentCreated, a); err != nil {
			return err
		}
		if a.Published() {
			return s.announce(ctx, tx, p, c, a)
		}
		return nil
	})
	return a, err
}

func (s AssignmentService) Get(ctx context.Context, p auth.Principal, id string) (AssignmentDetails, error) {
	var d AssignmentDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, c, st, err := s.visible(ctx, tx, p, id, false)
		if err != nil {
			return err
		}
		fileIDs, err := s.assignments.Files(ctx, tx, id)
		if err != nil {
			return err
		}
		attached, err := s.attachments.Load(ctx, tx, fileIDs)
		if err != nil {
			return err
		}
		d = AssignmentDetails{Assignment: a, CourseTitle: c.Title, Role: st.Shown(), Files: attached}
		if st.Teaches() {
			items := make([]AssignmentItem, 1)
			if err := s.fillCounts(ctx, tx, c.ID, []Assignment{a}, items); err != nil {
				return err
			}
			d.Counts = items[0].Counts
			return nil
		}
		people, err := s.people.ByIDs(ctx, tx, []string{p.UserID})
		if err != nil {
			return err
		}
		work, err := s.work(ctx, tx, a, people[0])
		d.Mine = &work
		return err
	})
	return d, err
}

func (s AssignmentService) Update(ctx context.Context, p auth.Principal, id string, in AssignmentInput) (Assignment, error) {
	in, err := in.normalized()
	if err != nil {
		return Assignment{}, err
	}
	var a Assignment
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		before, c, err := s.taught(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if before.Published() && !in.Published {
			return ErrInvalidAssignment
		}
		if before.Grading != in.Grading {
			graded, err := s.grades.ByAssignment(ctx, tx, []string{id}, nil)
			if err != nil {
				return err
			}
			if len(graded) > 0 {
				return ErrHasWork
			}
		}
		current, err := s.assignments.Files(ctx, tx, id)
		if err != nil {
			return err
		}
		if err := s.attachments.Check(ctx, tx, p.UserID, added(current, in.Files)); err != nil {
			return err
		}
		if a, err = s.assignments.Update(ctx, tx, id, in); err != nil {
			return err
		}
		if err := s.assignments.SetFiles(ctx, tx, id, in.Files); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionAssignmentUpdated, a); err != nil {
			return err
		}
		if !before.Published() && a.Published() {
			return s.announce(ctx, tx, p, c, a)
		}
		return nil
	})
	return a, err
}

func (s AssignmentService) Delete(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, _, err := s.taught(ctx, tx, p, id)
		if err != nil {
			return err
		}
		n, err := s.submissions.Count(ctx, tx, id)
		if err != nil {
			return err
		}
		if n > 0 {
			return ErrHasWork
		}
		if err := s.assignments.Delete(ctx, tx, id); err != nil {
			return err
		}
		return s.record(ctx, tx, p, ActionAssignmentDeleted, a)
	})
}

func (s AssignmentService) Submit(ctx context.Context, p auth.Principal, id string, in SubmissionInput) (Submission, bool, error) {
	in, err := in.normalized()
	if err != nil {
		return Submission{}, false, err
	}
	var sub Submission
	var created bool
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, c, st, err := s.visible(ctx, tx, p, id, true)
		if err != nil {
			return err
		}
		if st.Role != RoleStudent {
			return access.ErrForbidden
		}
		existing, err := s.submissions.ByClient(ctx, tx, id, p.UserID, in.ClientID)
		if err == nil {
			sub, err = s.withFiles(ctx, tx, existing)
			return err
		}
		if !errors.Is(err, ErrAssignmentNotFound) {
			return err
		}
		if err := s.teamless(ctx, tx, a, p.UserID); err != nil {
			return err
		}
		if err := s.acceptable(ctx, tx, p, a, c, in); err != nil {
			return err
		}
		help, err := s.helpFor(ctx, tx, a, p.UserID)
		if err != nil {
			return err
		}
		if sub, err = s.submissions.Insert(ctx, tx, id, p.UserID, in, nil, help); err != nil {
			return err
		}
		if err := s.submissions.SetFiles(ctx, tx, sub.ID, in.Files); err != nil {
			return err
		}
		if sub, err = s.withFiles(ctx, tx, sub); err != nil {
			return err
		}
		created = true
		name, err := s.displayName(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		return s.tellTeachers(ctx, tx, p.TenantID, c, a, name)
	})
	return sub, created, err
}

func (s AssignmentService) Work(ctx context.Context, p auth.Principal, id string) ([]StudentWork, error) {
	var list []StudentWork
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, c, err := s.taught(ctx, tx, p, id)
		if err != nil {
			return err
		}
		studentIDs, err := s.gate.StudentIDs(ctx, tx, c.ID)
		if err != nil {
			return err
		}
		people, err := s.people.ByIDs(ctx, tx, studentIDs)
		if err != nil {
			return err
		}
		latest, attempts, err := s.submissions.Latest(ctx, tx, id)
		if err != nil {
			return err
		}
		grades, err := s.grades.ByAssignment(ctx, tx, []string{a.ID}, nil)
		if err != nil {
			return err
		}
		byStudent := map[string]Grade{}
		for _, g := range grades {
			byStudent[g.StudentID] = g
		}
		list = make([]StudentWork, 0, len(people))
		for _, person := range people {
			w := StudentWork{Student: person, Attempts: attempts[person.UserID]}
			if sub, ok := latest[person.UserID]; ok {
				full, err := s.withFiles(ctx, tx, sub)
				if err != nil {
					return err
				}
				w.Latest = &full
			}
			if g, ok := byStudent[person.UserID]; ok {
				w.Grade = &g
			}
			list = append(list, w)
		}
		return nil
	})
	return list, err
}

func (s AssignmentService) Attempts(ctx context.Context, p auth.Principal, id, studentID string) ([]Submission, error) {
	var list []Submission
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		_, _, st, err := s.visible(ctx, tx, p, id, false)
		if err != nil {
			return err
		}
		if !st.Teaches() && studentID != p.UserID {
			return access.ErrForbidden
		}
		raw, err := s.submissions.Attempts(ctx, tx, id, studentID)
		if err != nil {
			return err
		}
		list = make([]Submission, 0, len(raw))
		for _, sub := range raw {
			full, err := s.withFiles(ctx, tx, sub)
			if err != nil {
				return err
			}
			list = append(list, full)
		}
		return nil
	})
	return list, err
}

func (s AssignmentService) Grade(ctx context.Context, p auth.Principal, id, studentID string, in GradeInput) (Grade, error) {
	var g Grade
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, c, err := s.taught(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if !a.Published() {
			return ErrNotPublished
		}
		if err := s.student(ctx, tx, c, studentID); err != nil {
			return err
		}
		g, err = s.grade(ctx, tx, p, a, c, studentID, in)
		return err
	})
	return g, err
}

func (s AssignmentService) grade(ctx context.Context, tx pgx.Tx, p auth.Principal, a Assignment, c Course, studentID string, in GradeInput) (Grade, error) {
	in, err := in.normalized(a)
	if err != nil {
		return Grade{}, err
	}
	g, err := s.grades.Upsert(ctx, tx, a.ID, studentID, p.UserID, in)
	if err != nil {
		return Grade{}, err
	}
	if err := s.log.Record(ctx, tx, p.UserID, audit.Change{
		Action: ActionGraded, TargetKind: "study.assignment", TargetID: a.ID,
		Details: map[string]any{"title": a.Title, "student": studentID},
	}); err != nil {
		return Grade{}, err
	}
	return g, s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
		Type: NotifyGraded, UserIDs: []string{studentID},
		Params: map[string]any{"assignment": a.Title, "course": c.Title},
		Link:   assignmentLink(a.ID),
	})
}

func (s AssignmentService) student(ctx context.Context, tx pgx.Tx, c Course, studentID string) error {
	students, err := s.gate.StudentIDs(ctx, tx, c.ID)
	if err != nil {
		return err
	}
	if !slices.Contains(students, studentID) {
		return ErrNotStudent
	}
	return nil
}

func (s AssignmentService) work(ctx context.Context, tx pgx.Tx, a Assignment, student identity.Person) (StudentWork, error) {
	w := StudentWork{Student: student}
	attempts, err := s.submissions.Attempts(ctx, tx, a.ID, student.UserID)
	if err != nil {
		return StudentWork{}, err
	}
	w.Attempts = len(attempts)
	if len(attempts) > 0 {
		latest, err := s.withFiles(ctx, tx, attempts[0])
		if err != nil {
			return StudentWork{}, err
		}
		w.Latest = &latest
	}
	grades, err := s.grades.ByAssignment(ctx, tx, []string{a.ID}, &student.UserID)
	if err != nil {
		return StudentWork{}, err
	}
	if len(grades) > 0 {
		w.Grade = &grades[0]
	}
	return w, nil
}

func (s AssignmentService) withFiles(ctx context.Context, tx pgx.Tx, sub Submission) (Submission, error) {
	ids, err := s.submissions.Files(ctx, tx, sub.ID)
	if err != nil {
		return Submission{}, err
	}
	sub.Files, err = s.attachments.Load(ctx, tx, ids)
	return sub, err
}

func (s AssignmentService) visible(ctx context.Context, tx pgx.Tx, p auth.Principal, id string, lock bool) (Assignment, Course, Standing, error) {
	a, err := s.assignments.ByID(ctx, tx, id, lock)
	if err != nil {
		return Assignment{}, Course{}, Standing{}, err
	}
	c, st, err := s.gate.Visible(ctx, tx, p, a.CourseID)
	if errors.Is(err, ErrCourseNotFound) || (err == nil && !a.Published() && !st.Teaches()) {
		return Assignment{}, Course{}, Standing{}, ErrAssignmentNotFound
	}
	if err != nil {
		return Assignment{}, Course{}, Standing{}, err
	}
	return a, c, st, nil
}

func (s AssignmentService) taught(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Assignment, Course, error) {
	a, c, st, err := s.visible(ctx, tx, p, id, true)
	if err != nil {
		return Assignment{}, Course{}, err
	}
	if !st.Teaches() {
		return Assignment{}, Course{}, access.ErrForbidden
	}
	return a, c, nil
}

func (s AssignmentService) announce(ctx context.Context, tx pgx.Tx, p auth.Principal, c Course, a Assignment) error {
	students, err := s.gate.StudentIDs(ctx, tx, c.ID)
	if err != nil {
		return err
	}
	if err := s.record(ctx, tx, p, ActionAssignmentPublished, a); err != nil {
		return err
	}
	return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
		Type: NotifyPublished, UserIDs: students,
		Params: map[string]any{"assignment": a.Title, "course": c.Title},
		Link:   assignmentLink(a.ID),
	})
}

func (s AssignmentService) tellTeachers(ctx context.Context, tx pgx.Tx, tenantID string, c Course, a Assignment, who string) error {
	teachers, err := s.gate.roster.Teachers(ctx, tx, c.ID)
	if err != nil {
		return err
	}
	return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
		Type: NotifySubmitted, UserIDs: teachers,
		Params: map[string]any{"assignment": a.Title, "course": c.Title, "student": who},
		Link:   assignmentLink(a.ID),
	})
}

func (s AssignmentService) displayName(ctx context.Context, tx pgx.Tx, userID string) (string, error) {
	people, err := s.people.ByIDs(ctx, tx, []string{userID})
	if err != nil || len(people) == 0 {
		return "", err
	}
	return people[0].DisplayName, nil
}

func (s AssignmentService) record(ctx context.Context, tx pgx.Tx, p auth.Principal, action string, a Assignment) error {
	return s.log.Record(ctx, tx, p.UserID, audit.Change{
		Action: action, TargetKind: "study.assignment", TargetID: a.ID,
		Details: map[string]any{"title": a.Title, "course": a.CourseID},
	})
}

func assignmentLink(id string) string {
	return "/study/assignments/" + id
}

func added(before, after []string) []string {
	var out []string
	for _, id := range after {
		if !slices.Contains(before, id) {
			out = append(out, id)
		}
	}
	return out
}
