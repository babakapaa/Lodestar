package league

type TransitionStatus string

const (
	TransitionPending  TransitionStatus = "pending"
	TransitionApplied  TransitionStatus = "applied"
	TransitionRejected TransitionStatus = "rejected"
	TransitionReverted TransitionStatus = "reverted"
)
