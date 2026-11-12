package module

type NotificationType struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	InApp    bool   `json:"inApp"`
	Push     bool   `json:"push"`
	Email    bool   `json:"email"`

	Names map[string]string           `json:"-"`
	Texts map[string]NotificationText `json:"-"`
}

type Notifying interface {
	NotificationTypes() []NotificationType
}
