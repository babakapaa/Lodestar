package setup_test

import (
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/setup"
)

func TestAssigneeIsInvitedAndRemindedOnce(t *testing.T) {
	w := newWorld(t)
	if err := w.setup.Begin(w.ctx, w.tenant); err != nil {
		t.Fatalf("begin: %v", err)
	}
	it := w.user()
	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly)
	later := time.Now().UTC().AddDate(0, 0, 10).Format(time.DateOnly)
	if _, err := w.setup.Assign(w.ctx, w.admin, setup.SectionUniversity, setup.AssignInput{UserID: &it.UserID, DueOn: &tomorrow}); err != nil {
		t.Fatalf("assign: %v", err)
	}
	if _, err := w.setup.Assign(w.ctx, w.admin, setup.SectionPeople, setup.AssignInput{UserID: &it.UserID, DueOn: &later}); err != nil {
		t.Fatalf("assign later: %v", err)
	}
	if _, err := w.setup.Assign(w.ctx, w.admin, setup.SectionPrivacy, setup.AssignInput{UserID: &w.admin.UserID, DueOn: &later}); err != nil {
		t.Fatalf("assign self: %v", err)
	}
	sent := w.sent.take()
	if len(sent) != 2 || sent[0].Type != setup.NotifyAssigned || sent[0].UserIDs[0] != it.UserID || sent[0].Link != "/admin/setup/university" {
		t.Fatalf("invitations = %+v", sent)
	}
	if _, err := w.setup.Assign(w.ctx, w.admin, setup.SectionUniversity, setup.AssignInput{UserID: &it.UserID, DueOn: &tomorrow}); err != nil {
		t.Fatalf("assign again: %v", err)
	}
	if sent := w.sent.take(); len(sent) != 0 {
		t.Fatalf("same assignee invited again = %+v", sent)
	}

	reminder := setup.NewReminder(w.pool, w.sent, time.Hour)
	for range 2 {
		if err := reminder.RemindAll(w.ctx); err != nil {
			t.Fatalf("remind: %v", err)
		}
	}
	sent = w.sent.take()
	if len(sent) != 1 || sent[0].Type != setup.NotifyDueSoon || sent[0].Params["due"] != tomorrow || sent[0].UserIDs[0] != it.UserID {
		t.Fatalf("reminders = %+v", sent)
	}
}
