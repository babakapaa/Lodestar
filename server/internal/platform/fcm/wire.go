package fcm

type request struct {
	Message message `json:"message"`
}

type message struct {
	Token        string            `json:"token"`
	Notification notification      `json:"notification"`
	Data         map[string]string `json:"data,omitempty"`
}

type notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type response struct {
	Error struct {
		Status  string `json:"status"`
		Details []struct {
			ErrorCode string `json:"errorCode"`
		} `json:"details"`
	} `json:"error"`
}

func (r response) unregistered() bool {
	for _, d := range r.Error.Details {
		if d.ErrorCode == "UNREGISTERED" {
			return true
		}
	}
	return false
}
