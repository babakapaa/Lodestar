package setup

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyAssigned = "setup.assigned"
	NotifyDueSoon  = "setup.due_soon"
)

func NotificationTypes() []module.NotificationType {
	return []module.NotificationType{
		{
			Key: NotifyAssigned, InApp: true, Push: true, Email: true,
			Names: map[string]string{"ru": "Поручения в мастере настройки", "kk": "Баптау шеберіндегі тапсырмалар", "en": "Setup wizard assignments"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Вам поручен раздел настройки Lodestar", Body: "Откройте мастер настройки, чтобы начать"},
				"kk": {Title: "Сізге Lodestar баптау бөлімі тапсырылды", Body: "Бастау үшін баптау шеберін ашыңыз"},
				"en": {Title: "You were assigned a Lodestar setup section", Body: "Open the setup wizard to start"},
			},
		},
		{
			Key: NotifyDueSoon, InApp: true, Push: true, Email: true,
			Names: map[string]string{"ru": "Сроки в мастере настройки", "kk": "Баптау шеберіндегі мерзімдер", "en": "Setup wizard due dates"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Срок по разделу настройки: {due}", Body: "Откройте мастер, чтобы закончить раздел"},
				"kk": {Title: "Баптау бөлімінің мерзімі: {due}", Body: "Бөлімді аяқтау үшін шеберді ашыңыз"},
				"en": {Title: "Setup section due {due}", Body: "Open the wizard to finish the section"},
			},
		},
	}
}
