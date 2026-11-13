package sandbox

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyTopupRequested = "sandbox.topup_requested"
	NotifyTopupApproved  = "sandbox.topup_approved"
	NotifyTopupDeclined  = "sandbox.topup_declined"
	NotifyRunFinished    = "sandbox.run_finished"
	NotifyRunFailed      = "sandbox.run_failed"
	NotifyIdle           = "sandbox.idle"
	NotifyQuotaLow       = "sandbox.quota_low"
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
		{
			Key: NotifyRunFinished, InApp: true, Push: true,
			Names: map[string]string{"ru": "Готовые задачи на GPU", "kk": "GPU-дағы дайын тапсырмалар", "en": "Finished GPU jobs"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Готово: {title}", Body: "Списано {hours} ч GPU"},
				"kk": {Title: "Дайын: {title}", Body: "{hours} сағ GPU есептен шығарылды"},
				"en": {Title: "Done: {title}", Body: "{hours} GPU hours used"},
			},
		},
		{
			Key: NotifyRunFailed, InApp: true, Push: true,
			Names: map[string]string{"ru": "Упавшие задачи на GPU", "kk": "GPU-дағы құлаған тапсырмалар", "en": "Failed GPU jobs"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Не получилось: {title}", Body: "Откройте запуск, чтобы узнать почему"},
				"kk": {Title: "Сәтсіз: {title}", Body: "Себебін білу үшін іске қосуды ашыңыз"},
				"en": {Title: "Failed: {title}", Body: "Open the run to see why"},
			},
		},
		{
			Key: NotifyIdle, InApp: true, Push: true,
			Names: map[string]string{"ru": "Остановка простаивающих сессий", "kk": "Бос тұрған сессияларды тоқтату", "en": "Idle sessions stopping"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "«{title}» остановится через {minutes} мин", Body: "Сессия простаивает. Файлы сохранятся"},
				"kk": {Title: "«{title}» {minutes} минуттан кейін тоқтайды", Body: "Сессия бос тұр. Файлдар сақталады"},
				"en": {Title: "“{title}” stops in {minutes} min", Body: "The session is idle. Your files are kept"},
			},
		},
		{
			Key: NotifyQuotaLow, InApp: true, Push: true,
			Names: map[string]string{"ru": "Заканчиваются GPU-часы", "kk": "GPU сағаттары таусылып барады", "en": "GPU hours running low"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Осталось меньше 10% GPU-часов", Body: "{hours} ч из {total}"},
				"kk": {Title: "GPU сағаттарының 10%-нан азы қалды", Body: "{total} сағаттың {hours} сағаты"},
				"en": {Title: "Less than 10% of GPU hours left", Body: "{hours} of {total} hours"},
			},
		},
	}
}
