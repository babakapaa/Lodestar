package study

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyPublished = "study.assignment_published"
	NotifySubmitted = "study.work_submitted"
	NotifyGraded    = "study.graded"
)

func (Module) NotificationTypes() []module.NotificationType {
	return append(workTypes(), scheduleTypes()...)
}

func workTypes() []module.NotificationType {
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

const (
	NotifyScheduleChanged = "study.schedule_changed"
	NotifyDeadline        = "study.deadline"
)

func scheduleTypes() []module.NotificationType {
	return []module.NotificationType{
		{
			Key: NotifyScheduleChanged, InApp: true, Push: true,
			Names: map[string]string{"ru": "Изменения в расписании", "kk": "Кестедегі өзгерістер", "en": "Schedule changes"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Изменение в расписании: {course}", Body: "{date}"},
				"kk": {Title: "Кестедегі өзгеріс: {course}", Body: "{date}"},
				"en": {Title: "Schedule change: {course}", Body: "{date}"},
			},
		},
		{
			Key: NotifyDeadline, InApp: true, Push: true,
			Names: map[string]string{"ru": "Скорые сроки сдачи", "kk": "Жақындаған тапсыру мерзімдері", "en": "Upcoming deadlines"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "Скоро срок: {assignment}", Body: "{course}"},
				"kk": {Title: "Мерзімі жақын: {assignment}", Body: "{course}"},
				"en": {Title: "Due soon: {assignment}", Body: "{course}"},
			},
		},
	}
}
