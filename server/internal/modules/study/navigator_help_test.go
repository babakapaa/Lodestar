package study_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

type fakeHelp struct {
	asked []string
}

func (f *fakeHelp) Help(_ context.Context, _ pgx.Tx, userID, scopeKind, scopeID string) (module.HelpSummary, error) {
	f.asked = append(f.asked, userID+"/"+scopeKind+"/"+scopeID)
	return module.HelpSummary{Replies: 3, Steps: module.HelpSteps{Question: 1, Source: 2}, Sources: []string{"Иванов И. И."}, Articles: []string{}}, nil
}

func TestNavigatorSummaryIsAttachedOnlyWhenLecturerAsks(t *testing.T) {
	w := newWorld(t)
	group := w.unit(nil, "group", "ИС-21")
	teacher := w.user("dana", "dana@kaznu.kz")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	w.join(group, student)
	c := w.course(nil, "Биология")
	w.add(study.ListUnits, c, group.ID)
	w.add(study.ListTeachers, c, teacher.UserID)
	help := &fakeHelp{}
	work := w.work.WithHelp(help)

	plain, err := work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Без сводки", Grading: study.GradingPass, Published: true})
	if err != nil || plain.NavigatorHelp {
		t.Fatalf("create plain = %+v, %v", plain, err)
	}
	sub, _, err := work.Submit(w.ctx, student, plain.ID, study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a03", Body: "Ответ"})
	if err != nil || sub.Help != nil || len(help.asked) != 0 {
		t.Fatalf("plain submit = %+v, %v, asked %v", sub, err, help.asked)
	}

	a, err := work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Со сводкой", Grading: study.GradingPass, Published: true, NavigatorHelp: true})
	if err != nil || !a.NavigatorHelp {
		t.Fatalf("create = %+v, %v", a, err)
	}
	sub, _, err = work.Submit(w.ctx, student, a.ID, study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a04", Body: "Ответ"})
	if err != nil || sub.Help == nil || sub.Help.Replies != 3 || sub.Help.Steps.Source != 2 {
		t.Fatalf("submit = %+v, %v", sub, err)
	}
	if len(help.asked) != 1 || help.asked[0] != student.UserID+"/assignment/"+a.ID {
		t.Fatalf("asked = %v", help.asked)
	}
	list, err := work.Work(w.ctx, teacher, a.ID)
	if err != nil || len(list) != 1 || list[0].Latest == nil || list[0].Latest.Help == nil || list[0].Latest.Help.Sources[0] != "Иванов И. И." {
		t.Fatalf("work = %+v, %v", list, err)
	}
}
