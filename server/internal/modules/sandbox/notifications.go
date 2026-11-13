package sandbox

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyTopupRequested = "sandbox.topup_requested"
	NotifyTopupApproved  = "sandbox.topup_approved"
	NotifyTopupDeclined  = "sandbox.topup_declined"
)

func (Module) NotificationTypes() []module.NotificationType {
	return []module.NotificationType{
		{
			Key: NotifyTopupRequested, InApp: true, Push: true,
			Names: map[string]string{"ru": "Запросы GPU-часов", "kk": "GPU сағаттарына сұраулар", "en": "GPU hour requests"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{student} просит {hours} ч GPU", Body: "{work}"},
				"kk": {Title: "{student} {hours} сағ GPU сұрайды", Body: "{work}"},
				"en": {Title: "{student} asks for {hours} GPU hours", Body: "{work}"},
			},
		},
		{
			Key: NotifyTopupApproved, InApp: true, Push: true,
			Names: map[string]string{"ru": "Добавленные GPU-часы", "kk": "Қосылған GPU сағаттары", "en": "GPU hours added"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Добавлено {hours} ч GPU", Body: "{work}"},
				"kk": {Title: "{hours} сағ GPU қосылды", Body: "{work}"},
				"en": {Title: "{hours} GPU hours added", Body: "{work}"},
			},
		},
		{
			Key: NotifyTopupDeclined, InApp: true, Push: true,
			Names: map[string]string{"ru": "Отклонённые запросы GPU-часов", "kk": "Қабылданбаған GPU сағаттары сұраулары", "en": "Declined GPU hour requests"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Запрос GPU-часов отклонён", Body: "{work}"},
				"kk": {Title: "GPU сағаттары сұрауы қабылданбады", Body: "{work}"},
				"en": {Title: "GPU hour request declined", Body: "{work}"},
			},
		},
	}
}
