package module

const (
	VoteAgreed   = "agreed"
	VoteObjected = "objected"
	VoteSilent   = "silent"
)

type TeamVote struct {
	UserID   string
	Decision string
	Comment  string
}
