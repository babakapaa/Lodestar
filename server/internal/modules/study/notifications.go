package study

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyPublished = "study.assignment_published"
	NotifySubmitted = "study.work_submitted"
	NotifyGraded    = "study.graded"
)

func (Module) NotificationTypes() []module.NotificationType {
	return []module.NotificationType{
		{
			Key: NotifyPublished, InApp: true, Push: true,
			Names: map[string]string{"ru": "Новые задания", "kk": "Жаңа тапсырмалар", "en": "New assignments"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Новое задание: {assignment}", Body: "{course}"},
				"kk": {Title: "Жаңа тапсырма: {assignment}", Body: "{course}"},
				"en": {Title: "New assignment: {assignment}", Body: "{course}"},
			},
		},
		{
			Key: NotifySubmitted, InApp: true,
			Names: map[string]string{"ru": "Сданные работы", "kk": "Тапсырылған жұмыстар", "en": "Submitted work"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Сдана работа: {assignment}", Body: "{student}, {course}"},
				"kk": {Title: "Жұмыс тапсырылды: {assignment}", Body: "{student}, {course}"},
				"en": {Title: "Work submitted: {assignment}", Body: "{student}, {course}"},
			},
		},
		{
			Key: NotifyGraded, InApp: true, Push: true,
			Names: map[string]string{"ru": "Оценки и отзывы", "kk": "Бағалар мен пікірлер", "en": "Grades and feedback"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Работа проверена: {assignment}", Body: "{course}"},
				"kk": {Title: "Жұмыс тексерілді: {assignment}", Body: "{course}"},
				"en": {Title: "Work reviewed: {assignment}", Body: "{course}"},
			},
		},
	}
}
