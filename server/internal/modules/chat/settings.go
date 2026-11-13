package chat

import "github.com/babakapaa/lodestar/server/internal/core/schema"

const (
	RetentionSetting   = "chat.retention_days"
	CallsSetting       = "chat.calls"
	ScreenShareSetting = "chat.screen_share"
	RecordingSetting   = "chat.call_recording"
	ConsentSetting     = "chat.recording_consent"
)

func (Module) Settings() []schema.Definition {
	minDays, maxDays := 0, 36500
	return []schema.Definition{
		{Key: RetentionSetting, Type: schema.Integer, Default: 0, Min: &minDays, Max: &maxDays},
		{Key: CallsSetting, Type: schema.Bool, Default: true},
		{Key: ScreenShareSetting, Type: schema.Choice, Options: []string{ShareEveryone, ShareAdmins}, Default: ShareEveryone},
		{Key: RecordingSetting, Type: schema.Choice, Options: []string{RecordOff, RecordAdmins, RecordEveryone}, Default: RecordOff},
		{Key: ConsentSetting, Type: schema.Bool, Default: true},
	}
}
