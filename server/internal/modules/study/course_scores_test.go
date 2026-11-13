package study_test

import (
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestCourseScoresCountOnlyDefendedWork(t *testing.T) {
	d := newDefenseCase(t, study.DefenseFull)
	w := d.w
	quiz, err := w.work.Create(w.ctx, d.teacher, d.course.ID, study.AssignmentInput{Title: "Тест", Grading: study.GradingPass, Published: true})
	if err != nil {
		t.Fatalf("quiz: %v", err)
	}
	if _, err := w.work.Create(w.ctx, d.teacher, d.course.ID, study.AssignmentInput{Title: "Черновик", Grading: study.GradingPass}); err != nil {
		t.Fatalf("draft: %v", err)
	}
	passed := true
	if _, err := w.work.Grade(w.ctx, d.teacher, quiz.ID, d.first.UserID, study.GradeInput{Passed: &passed}); err != nil {
		t.Fatalf("grade quiz: %v", err)
	}
	early := 90.0
	if _, err := w.work.Grade(w.ctx, d.teacher, d.work.ID, d.second.UserID, study.GradeInput{Points: &early}); err != nil {
		t.Fatalf("grade before defense: %v", err)
	}
	points := 80.0
	if _, err := d.defenses.Finish(w.ctx, d.teacher, d.work.ID, d.first.UserID, study.FinishInput{
		Status: study.DefenseHeld, Grade: &study.GradeInput{Points: &points},
	}); err != nil {
		t.Fatalf("finish: %v", err)
	}

	var works []module.ScoredWork
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		works, err = study.CourseScores{}.ScoredWorks(w.ctx, tx, d.course.ID)
		return err
	})
	if err != nil || len(works) != 2 {
		t.Fatalf("works = %+v, %v", works, err)
	}
	byID := map[string]module.ScoredWork{}
	for _, work := range works {
		byID[work.ID] = work
	}
	project := byID[d.work.ID]
	if !project.Project || len(project.Results) != 1 || project.Results[0].UserID != d.first.UserID || project.Results[0].Share() != 0.8 {
		t.Fatalf("project = %+v", project)
	}
	test := byID[quiz.ID]
	if test.Project || len(test.Results) != 1 || test.Results[0].Share() != 1 {
		t.Fatalf("quiz = %+v", test)
	}
}
