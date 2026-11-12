package study

import "time"

type SearchMapDraft struct {
	Map       SearchMap  `json:"map"`
	Rules     *MapRules  `json:"rules"`
	Missing   []string   `json:"missing"`
	ChangedAt *time.Time `json:"changedAt"`
}
