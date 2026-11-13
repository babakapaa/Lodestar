package chat

import "time"

func recordingFile(now time.Time, audioOnly bool) (string, string) {
	name := "call-" + now.UTC().Format("2006-01-02-1504")
	if audioOnly {
		return name + ".ogg", "audio/ogg"
	}
	return name + ".mp4", "video/mp4"
}
