package livekit

func fileType(audioOnly bool) string {
	if audioOnly {
		return "OGG"
	}
	return "MP4"
}
