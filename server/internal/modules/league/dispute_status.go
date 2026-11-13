package league

type DisputeStatus string

const (
	DisputeOpen     DisputeStatus = "open"
	DisputeAccepted DisputeStatus = "accepted"
	DisputeRejected DisputeStatus = "rejected"
)
