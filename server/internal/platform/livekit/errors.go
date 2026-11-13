package livekit

import "errors"

var ErrBadWebhook = errors.New("webhook signature does not match")
