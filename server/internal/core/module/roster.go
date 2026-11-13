package module

type Roster struct {
	Kind     string
	ID       string
	Title    string
	Leads    []string
	Members  []string
	Group    bool
	Archived bool
}
