package audit

type Change struct {
	Action     string
	TargetKind string
	TargetID   string
	Details    any
}
