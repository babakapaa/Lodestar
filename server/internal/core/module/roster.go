package module

type Roster struct {
	Kind     string
	ID       string
	Title    string
	Leads    []string
	Members  []string
	Archived bool
}
