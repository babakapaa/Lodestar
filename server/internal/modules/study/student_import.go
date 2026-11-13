package study

type StudentImport struct {
	Added    int           `json:"added"`
	Existing int           `json:"existing"`
	Errors   []ImportError `json:"errors"`
}

type ImportError struct {
	Line int    `json:"line"`
	Code string `json:"code"`
}

const (
	ImportMissingEmail  = "missing_email"
	ImportUnknownUser   = "unknown_user"
	ImportAmbiguousUser = "ambiguous_user"
)
