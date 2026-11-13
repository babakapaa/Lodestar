package chat

import "github.com/babakapaa/lodestar/server/internal/core/module"

const NotifyMessage = "chat.message"

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
	}
}
