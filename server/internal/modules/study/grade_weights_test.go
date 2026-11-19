package study_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestGradeIsAddedUpFromWeightedParts(t *testing.T) {
	w := newWorld(t)
	group := w.unit(nil, "group", "ИС-21")
	teacher := w.user("dana", "dana@kaznu.kz")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	w.join(group, student)
	c := w.course(nil, "Биология")
	w.add(study.ListUnits, c, group.ID)
	w.add(study.ListTeachers, c, teacher.UserID)
	ten := 10.0

	if _, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Зачёт", Grading: study.GradingPass, Weights: &study.GradeWeights{Defense: 50, Map: 30, Work: 20}}); !errors.Is(err, study.ErrInvalidAssignment) {
		t.Fatalf("weights on pass grading = %v", err)
	}
	if _, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Эссе", Grading: study.GradingPoints, MaxPoints: &ten, Weights: &study.GradeWeights{Defense: 50, Map: 30, Work: 10}}); !errors.Is(err, study.ErrInvalidAssignment) {
		t.Fatalf("weights not adding up = %v", err)
	}
	a, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Эссе", Grading: study.GradingPoints, MaxPoints: &ten, Weights: &study.GradeWeights{Defense: 50, Map: 30, Work: 20}, Published: true})
	if err != nil || a.Weights == nil || a.Weights.Defense != 50 {
		t.Fatalf("create = %+v, %v", a, err)
	}

	defense, mapPart, work := 80.0, 100.0, 50.0
	if _, err := w.work.Grade(w.ctx, teacher, a.ID, student.UserID, study.GradeInput{Parts: &study.GradeParts{Defense: &defense, Map: &mapPart}}); !errors.Is(err, study.ErrInvalidGrade) {
		t.Fatalf("missing part = %v", err)
	}
	g, err := w.work.Grade(w.ctx, teacher, a.ID, student.UserID, study.GradeInput{Parts: &study.GradeParts{Defense: &defense, Map: &mapPart, Work: &work}})
	if err != nil || g.Points == nil || *g.Points != 8 || g.Parts == nil || *g.Parts.Work != 50 {
		t.Fatalf("grade = %+v, %v", g, err)
	}
	seven := 7.0
	g, err = w.work.Grade(w.ctx, teacher, a.ID, student.UserID, study.GradeInput{Points: &seven})
	if err != nil || *g.Points != 7 || g.Parts != nil {
		t.Fatalf("lecturer's own mark = %+v, %v", g, err)
	}
}
