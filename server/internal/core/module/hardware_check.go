package module

type HardwareCheck struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Status   string `json:"status"`
	Detail   string `json:"detail"`
	Millis   int64  `json:"millis"`
}

const (
	CheckOK     = "ok"
	CheckOff    = "off"
	CheckFailed = "failed"
)
