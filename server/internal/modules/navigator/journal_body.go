package navigator

type journalBody struct {
	Query   string          `json:"query"`
	Sources []JournalSource `json:"sources,omitempty"`
}
