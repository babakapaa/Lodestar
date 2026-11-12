package structure

type ImportReport struct {
	Created int           `json:"created"`
	Updated int           `json:"updated"`
	Errors  []ImportError `json:"errors"`
}

type ImportError struct {
	Line int    `json:"line"`
	Code string `json:"code"`
}

const (
	ImportMissingParent = "missing_parent"
	ImportInvalidUnit   = "invalid_unit"
	ImportCycle         = "cycle"
	ImportMissingUnit   = "missing_unit"
	ImportUnknownUser   = "unknown_user"
	ImportAmbiguousUser = "ambiguous_user"
)
