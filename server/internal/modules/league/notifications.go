package league

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyPromoted  = "league.promoted"
	NotifyRelegated = "league.relegated"
	NotifyQueued    = "league.queued"
	NotifyDisputed  = "league.disputed"
	NotifyAnswered  = "league.answered"
)

func (Module) NotificationTypes() []module.NotificationType {
	return []module.NotificationType{
		{
			Key: NotifyPromoted, InApp: true, Push: true,
			Names: map[string]string{"ru": "Переход в дивизион выше", "kk": "Жоғары дивизионға өту", "en": "Moving up a division"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{league}", Body: "Вы перешли в дивизион {division}"},
				"kk": {Title: "{league}", Body: "Сіз {division} дивизионына өттіңіз"},
				"en": {Title: "{league}", Body: "You moved up to division {division}"},
			},
		},
		{
			Key: NotifyRelegated, InApp: true, Push: true,
			Names: map[string]string{"ru": "Переход в дивизион ниже", "kk": "Төмен дивизионға өту", "en": "Moving down a division"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{league}", Body: "По итогам сезона вы в дивизионе {division}"},
				"kk": {Title: "{league}", Body: "Маусым қорытындысы бойынша сіз {division} дивизионындасыз"},
				"en": {Title: "{league}", Body: "After the season you are in division {division}"},
			},
		},
		{
			Key: NotifyQueued, InApp: true,
			Names: map[string]string{"ru": "Переходы на подтверждение", "kk": "Растауды күтетін ауысулар", "en": "Moves to confirm"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{league}", Body: "Есть переходы между дивизионами, которые ждут вашего решения"},
				"kk": {Title: "{league}", Body: "Дивизиондар арасындағы ауысулар сіздің шешіміңізді күтуде"},
				"en": {Title: "{league}", Body: "Division moves are waiting for your decision"},
			},
		},
		{
			Key: NotifyDisputed, InApp: true, Push: true,
			Names: map[string]string{"ru": "Споры о баллах лиги", "kk": "Лига ұпайлары бойынша даулар", "en": "League point disputes"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{league}", Body: "{student} оспаривает баллы или переход"},
				"kk": {Title: "{league}", Body: "{student} ұпайды немесе ауысуды даулайды"},
				"en": {Title: "{league}", Body: "{student} disputes points or a move"},
			},
		},
		{
			Key: NotifyAnswered, InApp: true, Push: true,
			Names: map[string]string{"ru": "Ответы на споры в лиге", "kk": "Лигадағы дауларға жауаптар", "en": "Answers to league disputes"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{league}", Body: "Преподаватель ответил на ваш спор"},
				"kk": {Title: "{league}", Body: "Оқытушы дауыңызға жауап берді"},
				"en": {Title: "{league}", Body: "Your teacher answered your dispute"},
			},
		},
	}
}
