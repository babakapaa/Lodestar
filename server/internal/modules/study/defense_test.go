package study_test

import (
	"errors"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

type defenseCase struct {
	w        *world
	defenses study.DefenseService
	teacher  auth.Principal
	first    auth.Principal
	second   auth.Principal
	course   study.Course
	work     study.Assignment
}

func newDefenseCase(t *testing.T, format study.DefenseFormat) defenseCase {
	t.Helper()
	w := newWorld(t)
	group := w.unit(nil, "group", "ИС-21")
	teacher := w.user("dana", "dana@kaznu.kz")
	first := w.user("aigerim", "aigerim@kaznu.kz")
	second := w.user("erlan", "erlan@kaznu.kz")
	w.join(group, first)
	w.join(group, second)
	c := w.course(nil, "Базы данных")
	w.add(study.ListUnits, c, group.ID)
	w.add(study.ListTeachers, c, teacher.UserID)
	maxPoints := 100.0
	a, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{
		Title: "Курсовая", Grading: study.GradingPoints, MaxPoints: &maxPoints, Defense: &format, Published: true,
	})
	if err != nil {
		t.Fatalf("assignment: %v", err)
	}
	return defenseCase{w: w, defenses: study.NewDefenseService(w.kernel()), teacher: teacher, first: first, second: second, course: c, work: a}
}

func TestStudentsBookDefenseAndTeacherKeepsProtocolWithGrade(t *testing.T) {
	d := newDefenseCase(t, study.DefenseFull)
	w := d.w
	start := time.Now().Add(24 * time.Hour).Truncate(time.Minute)
	slots, err := d.defenses.AddSlots(w.ctx, d.teacher, d.work.ID, study.SlotsInput{StartsAt: start, Minutes: 15, Place: "Ауд. 214", Count: 2})
	if err != nil || len(slots) != 2 || !slots[1].StartsAt.Equal(start.Add(15*time.Minute)) {
		t.Fatalf("slots = %+v, %v", slots, err)
	}
	if _, err := d.defenses.AddSlots(w.ctx, d.first, d.work.ID, study.SlotsInput{StartsAt: start, Minutes: 15}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student adds slots: %v", err)
	}

	my, err := d.defenses.Book(w.ctx, d.first, d.work.ID, study.BookingInput{SlotID: slots[0].ID})
	if err != nil || my.Slot == nil || my.Slot.ID != slots[0].ID || !my.Required || my.Recorded {
		t.Fatalf("book = %+v, %v", my, err)
	}
	if _, err := d.defenses.Book(w.ctx, d.second, d.work.ID, study.BookingInput{SlotID: slots[0].ID}); !errors.Is(err, study.ErrSlotFull) {
		t.Fatalf("second in full slot: %v", err)
	}
	if my, err = d.defenses.Book(w.ctx, d.first, d.work.ID, study.BookingInput{SlotID: slots[1].ID}); err != nil || my.Slot.ID != slots[1].ID {
		t.Fatalf("move = %+v, %v", my, err)
	}
	if _, err := d.defenses.Book(w.ctx, d.second, d.work.ID, study.BookingInput{SlotID: slots[0].ID}); err != nil {
		t.Fatalf("second books freed slot: %v", err)
	}

	seen, err := d.defenses.Slots(w.ctx, d.second, d.work.ID)
	if err != nil || !seen[0].Mine || seen[1].Mine || seen[1].Taken != 1 || seen[0].Students != nil {
		t.Fatalf("student slots = %+v, %v", seen, err)
	}
	taught, err := d.defenses.Slots(w.ctx, d.teacher, d.work.ID)
	if err != nil || len(taught[1].Students) != 1 || taught[1].Students[0].UserID != d.first.UserID {
		t.Fatalf("teacher slots = %+v, %v", taught, err)
	}

	mark := 2
	saved, err := d.defenses.Save(w.ctx, d.teacher, d.work.ID, d.first.UserID, study.DefenseInput{
		Questions: []study.DefenseQuestion{
			{Kind: study.QuestionPath, Text: "Почему вы отбросили статью Кодда 1970 года?"},
			{Kind: study.QuestionUnderstanding, Text: "Что такое третья нормальная форма?", Mark: &mark, Note: "Ответила уверенно"},
		},
	})
	if err != nil || len(saved.Questions) != 2 {
		t.Fatalf("save = %+v, %v", saved, err)
	}
	if my, _ := d.defenses.Mine(w.ctx, d.first, d.work.ID); my.Protocol != nil {
		t.Fatalf("questions leaked before defense: %+v", my.Protocol)
	}
	if _, err := d.defenses.Get(w.ctx, d.first, d.work.ID, d.first.UserID); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student reads draft questions: %v", err)
	}

	points := 87.0
	held, err := d.defenses.Finish(w.ctx, d.teacher, d.work.ID, d.first.UserID, study.FinishInput{
		Status: study.DefenseHeld, Grade: &study.GradeInput{Points: &points, Comment: "Защита уверенная"},
	})
	if err != nil || held.Status != study.DefenseHeld || held.HeldAt == nil {
		t.Fatalf("finish = %+v, %v", held, err)
	}
	my, err = d.defenses.Mine(w.ctx, d.first, d.work.ID)
	if err != nil || my.Protocol == nil || len(my.Protocol.Questions) != 2 || my.Status != study.DefenseHeld {
		t.Fatalf("protocol = %+v, %v", my, err)
	}
	if _, err := d.defenses.Unbook(w.ctx, d.first, d.work.ID); !errors.Is(err, study.ErrDefenseFinished) {
		t.Fatalf("unbook after defense: %v", err)
	}
	if n := len(w.sent.of(study.NotifyDefenseHeld)); n != 1 {
		t.Fatalf("held notified %d times", n)
	}

	rows, err := d.defenses.List(w.ctx, d.teacher, d.work.ID)
	if err != nil || len(rows) != 2 {
		t.Fatalf("rows = %+v, %v", rows, err)
	}
	for _, row := range rows {
		if row.Student.UserID == d.first.UserID && (row.Grade == nil || *row.Grade.Points != 87 || row.Status != study.DefenseHeld) {
			t.Fatalf("first row = %+v", row)
		}
	}

	if err := d.defenses.DeleteSlot(w.ctx, d.teacher, slots[0].ID); err != nil {
		t.Fatalf("delete slot: %v", err)
	}
	cancelled := w.sent.of(study.NotifyDefenseCancelled)
	if len(cancelled) != 1 || cancelled[0].UserIDs[0] != d.second.UserID {
		t.Fatalf("cancelled = %+v", cancelled)
	}
	if my, _ := d.defenses.Mine(w.ctx, d.second, d.work.ID); my.Slot != nil {
		t.Fatalf("slot kept after delete: %+v", my.Slot)
	}
}

func TestSampleDefenseOnlyForChosenStudents(t *testing.T) {
	d := newDefenseCase(t, study.DefenseSample)
	w := d.w
	slots, err := d.defenses.AddSlots(w.ctx, d.teacher, d.work.ID, study.SlotsInput{StartsAt: time.Now().Add(time.Hour), Minutes: 20, Capacity: 2})
	if err != nil {
		t.Fatalf("slots: %v", err)
	}
	if my, err := d.defenses.Mine(w.ctx, d.first, d.work.ID); err != nil || my.Required {
		t.Fatalf("required before sample: %+v, %v", my, err)
	}
	if _, err := d.defenses.Book(w.ctx, d.first, d.work.ID, study.BookingInput{SlotID: slots[0].ID}); !errors.Is(err, study.ErrNotChosen) {
		t.Fatalf("unchosen books: %v", err)
	}
	chosen, err := d.defenses.Sample(w.ctx, d.teacher, d.work.ID, study.SampleInput{Count: 5})
	if err != nil || len(chosen) != 2 {
		t.Fatalf("sample = %v, %v", chosen, err)
	}
	if again, _ := d.defenses.Sample(w.ctx, d.teacher, d.work.ID, study.SampleInput{Count: 1}); len(again) != 0 {
		t.Fatalf("sampled twice: %v", again)
	}
	if _, err := d.defenses.Book(w.ctx, d.first, d.work.ID, study.BookingInput{SlotID: slots[0].ID}); err != nil {
		t.Fatalf("chosen books: %v", err)
	}
	if _, err := d.defenses.Select(w.ctx, d.teacher, d.work.ID, d.second.UserID, study.SelectionInput{Selected: false}); err != nil {
		t.Fatalf("unselect: %v", err)
	}
	if my, _ := d.defenses.Mine(w.ctx, d.second, d.work.ID); my.Required {
		t.Fatal("unselected still required")
	}
}

func TestDefenseSlotsAreClosedOnceStarted(t *testing.T) {
	d := newDefenseCase(t, study.DefenseShort)
	w := d.w
	slots, err := d.defenses.AddSlots(w.ctx, d.teacher, d.work.ID, study.SlotsInput{StartsAt: time.Now().Add(-time.Minute), Minutes: 10})
	if err != nil {
		t.Fatalf("slots: %v", err)
	}
	if _, err := d.defenses.Book(w.ctx, d.first, d.work.ID, study.BookingInput{SlotID: slots[0].ID}); !errors.Is(err, study.ErrSlotStarted) {
		t.Fatalf("book started slot: %v", err)
	}
	if _, err := d.defenses.AddSlots(w.ctx, d.teacher, d.work.ID, study.SlotsInput{StartsAt: time.Now(), Minutes: 1}); !errors.Is(err, study.ErrInvalidSlot) {
		t.Fatalf("short slot: %v", err)
	}
	if _, err := d.defenses.Finish(w.ctx, d.teacher, d.work.ID, d.first.UserID, study.FinishInput{Status: study.DefensePlanned}); !errors.Is(err, study.ErrInvalidDefense) {
		t.Fatalf("finish as planned: %v", err)
	}
	none, err := d.w.work.Create(w.ctx, d.teacher, d.course.ID, study.AssignmentInput{Title: "Без защиты", Published: true})
	if err != nil {
		t.Fatalf("plain assignment: %v", err)
	}
	if _, err := d.defenses.Slots(w.ctx, d.first, none.ID); !errors.Is(err, study.ErrNoDefense) {
		t.Fatalf("slots without defense: %v", err)
	}
}
