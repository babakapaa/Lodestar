package setup

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type Overview struct {
	Launched        bool             `json:"launched"`
	LaunchedAt      *time.Time       `json:"launchedAt"`
	ProtocolVersion int              `json:"protocolVersion"`
	SignedBy        *string          `json:"signedBy"`
	SignedAt        *time.Time       `json:"signedAt"`
	Sections        []SectionState   `json:"sections"`
	Legal           []LegalAnswer    `json:"legal"`
	Problems        []module.Problem `json:"problems"`
}
