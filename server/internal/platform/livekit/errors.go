package livekit

import "errors"

var (
	ErrBadWebhook       = errors.New("webhook signature does not match")
	ErrNoRecording      = errors.New("recording storage is not configured")
	ErrRecordingRefused = errors.New("livekit refused the recording")
)
