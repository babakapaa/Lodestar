package team

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyInvited        = "team.invited"
	NotifyLeaveRequested = "team.leave_requested"
	NotifyLeaveApproved  = "team.leave_approved"
	NotifyLeaveRejected  = "team.leave_rejected"
	NotifyChanged        = "team.changed"
)

func (Module) NotificationTypes() []module.NotificationType {
	return []module.NotificationType{
		{
			Key: NotifyInvited, InApp: true, Push: true,
			Names: map[string]string{"ru": "Приглашения в команду", "kk": "Командаға шақырулар", "en": "Team invitations"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{scope}", Body: "{inviter} зовёт вас в команду «{team}»"},
				"kk": {Title: "{scope}", Body: "{inviter} сізді «{team}» командасына шақырады"},
				"en": {Title: "{scope}", Body: "{inviter} invites you to team “{team}”"},
			},
		},
		{
			Key: NotifyLeaveRequested, InApp: true, Push: true,
			Names: map[string]string{"ru": "Просьбы выйти из команды", "kk": "Командадан шығу өтініштері", "en": "Requests to leave a team"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{scope}", Body: "{student} просит выйти из команды «{team}»"},
				"kk": {Title: "{scope}", Body: "{student} «{team}» командасынан шығуды сұрайды"},
				"en": {Title: "{scope}", Body: "{student} asks to leave team “{team}”"},
			},
		},
		{
			Key: NotifyLeaveApproved, InApp: true, Push: true,
			Names: map[string]string{"ru": "Решения о выходе из команды", "kk": "Командадан шығу шешімдері", "en": "Decisions on leaving a team"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{scope}", Body: "Преподаватель отпустил вас из команды «{team}»"},
				"kk": {Title: "{scope}", Body: "Оқытушы сізді «{team}» командасынан босатты"},
				"en": {Title: "{scope}", Body: "Your teacher let you leave team “{team}”"},
			},
		},
		{
			Key: NotifyLeaveRejected, InApp: true, Push: true,
			Names: map[string]string{"ru": "Отказы в выходе из команды", "kk": "Командадан шығудан бас тарту", "en": "Refused team leave requests"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{scope}", Body: "Преподаватель оставил вас в команде «{team}»"},
				"kk": {Title: "{scope}", Body: "Оқытушы сізді «{team}» командасында қалдырды"},
				"en": {Title: "{scope}", Body: "Your teacher kept you in team “{team}”"},
			},
		},
		{
			Key: NotifyChanged, InApp: true, Push: true,
			Names: map[string]string{"ru": "Изменения состава команды", "kk": "Команда құрамының өзгерісі", "en": "Team changes"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{scope}", Body: "Преподаватель изменил ваш состав: «{team}»"},
				"kk": {Title: "{scope}", Body: "Оқытушы құрамыңызды өзгертті: «{team}»"},
				"en": {Title: "{scope}", Body: "Your teacher changed your team: “{team}”"},
			},
		},
	}
}

const invitationsLink = "/teams"

func teamLink(id string) string {
	return "/teams/" + id
}
