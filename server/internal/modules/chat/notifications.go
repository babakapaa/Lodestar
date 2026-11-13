package chat

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	NotifyMessage = "chat.message"
	NotifyMention = "chat.mention"
)

func (Module) NotificationTypes() []module.NotificationType {
	return []module.NotificationType{
		{
			Key: NotifyMessage, Push: true,
			Names: map[string]string{"ru": "Сообщения в чатах", "kk": "Чаттардағы хабарламалар", "en": "Chat messages"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{chat}", Body: "Новое сообщение от {author}"},
				"kk": {Title: "{chat}", Body: "{author} жаңа хабарлама жіберді"},
				"en": {Title: "{chat}", Body: "New message from {author}"},
			},
		},
		{
			Key: NotifyMention, InApp: true, Push: true,
			Names: map[string]string{"ru": "Упоминания в чатах", "kk": "Чаттардағы атаулар", "en": "Chat mentions"},
			Texts: map[string]module.NotificationText{
				"ru": {Title: "{chat}", Body: "{author} упомянул(а) вас"},
				"kk": {Title: "{chat}", Body: "{author} сізді атап өтті"},
				"en": {Title: "{chat}", Body: "{author} mentioned you"},
			},
		},
	}
}
