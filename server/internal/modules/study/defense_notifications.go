package study

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyDefenseCancelled = "study.defense_cancelled"
	NotifyDefenseHeld      = "study.defense_held"
)

func defenseTypes() []module.NotificationType {
	return []module.NotificationType{
		{
			Key: NotifyDefenseCancelled, InApp: true, Push: true,
			Names: map[string]string{"ru": "Отменённые защиты", "kk": "Қорғаудың күшін жою", "en": "Cancelled defenses"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Время защиты отменено: {assignment}", Body: "Запишитесь на другое время. {course}"},
				"kk": {Title: "Қорғау уақыты жойылды: {assignment}", Body: "Басқа уақытқа жазылыңыз. {course}"},
				"en": {Title: "Defense time cancelled: {assignment}", Body: "Please book another time. {course}"},
			},
		},
		{
			Key: NotifyDefenseHeld, InApp: true, Push: true,
			Names: map[string]string{"ru": "Протоколы защит", "kk": "Қорғау хаттамалары", "en": "Defense reports"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Протокол защиты готов: {assignment}", Body: "{course}"},
				"kk": {Title: "Қорғау хаттамасы дайын: {assignment}", Body: "{course}"},
				"en": {Title: "Defense report is ready: {assignment}", Body: "{course}"},
			},
		},
	}
}
