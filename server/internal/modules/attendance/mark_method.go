package attendance

type MarkMethod string

const (
	ByQR     MarkMethod = "qr"
	ByCode   MarkMethod = "code"
	ByList   MarkMethod = "list"
	ByManual MarkMethod = "manual"
)
