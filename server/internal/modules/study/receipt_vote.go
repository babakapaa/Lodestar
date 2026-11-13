package study

import "github.com/babakapaa/lodestar/server/internal/core/identity"

type ReceiptVote struct {
	identity.Person
	Decision string `json:"decision"`
	Comment  string `json:"comment"`
}
