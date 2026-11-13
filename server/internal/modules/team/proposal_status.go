package team

type ProposalStatus string

const (
	ProposalOpen     ProposalStatus = "open"
	ProposalSent     ProposalStatus = "sent"
	ProposalReplaced ProposalStatus = "replaced"
)
