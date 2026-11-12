package notes

type VersionContent struct {
	Version
	Data []byte `json:"data"`
	Body string `json:"body"`
}
