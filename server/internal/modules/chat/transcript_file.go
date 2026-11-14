package chat

import (
	"path"
	"strings"
)

func transcriptFile(recording string) string {
	return strings.TrimSuffix(recording, path.Ext(recording)) + "-transcript.txt"
}
