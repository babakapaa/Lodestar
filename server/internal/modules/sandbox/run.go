package sandbox

import "time"

type Run struct {
	ID           string     `json:"id"`
	UserID       string     `json:"userId"`
	DisplayName  string     `json:"displayName"`
	TeamID       *string    `json:"teamId"`
	TeamName     string     `json:"teamName"`
	Kind         string     `json:"kind"`
	Title        string     `json:"title"`
	Command      string     `json:"command"`
	TemplateID   string     `json:"templateId"`
	TemplateName string     `json:"templateName"`
	Image        string     `json:"-"`
	ClassID      string     `json:"classId"`
	ClassName    string     `json:"className"`
	Rate         float64    `json:"rate"`
	MaxHours     float64    `json:"maxHours"`
	WorkKind     *string    `json:"workKind"`
	WorkID       *string    `json:"workId"`
	WorkTitle    string     `json:"workTitle"`
	DueAt        *time.Time `json:"dueAt"`
	Status       string     `json:"status"`
	Position     *int       `json:"position"`
	URL          string     `json:"url"`
	ExitCode     *int       `json:"exitCode"`
	Reason       string     `json:"reason"`
	Charged      *float64   `json:"charged"`
	QueuedAt     time.Time  `json:"queuedAt"`
	StartedAt    *time.Time `json:"startedAt"`
	EndedAt      *time.Time `json:"endedAt"`
	HeldAt       *time.Time `json:"heldAt"`
	HoldNote     string     `json:"holdNote"`
	Refunded     bool       `json:"refunded"`
	GPULoad      *float64   `json:"gpuLoad"`
	Log          string     `json:"log,omitempty"`
	ActiveAt     *time.Time `json:"-"`
	NetBytes     int64      `json:"-"`
	IdleWarned   bool       `json:"-"`
	HeldSecs     float64    `json:"-"`
	NodeID       *string    `json:"-"`
	TermID       string     `json:"-"`
}

func (r Run) owner() (string, string) {
	if r.TeamID != nil {
		return OwnerTeam, *r.TeamID
	}
	return OwnerUser, r.UserID
}

func (r Run) reserve() float64 {
	return hours(r.MaxHours * r.Rate)
}

func (r Run) used(end time.Time) time.Duration {
	if r.StartedAt == nil {
		return 0
	}
	if r.HeldAt != nil && end.After(*r.HeldAt) {
		end = *r.HeldAt
	}
	return max(end.Sub(*r.StartedAt)-time.Duration(r.HeldSecs*float64(time.Second)), 0)
}

func (r Run) cost(end time.Time, t tariff) float64 {
	if r.StartedAt == nil {
		return 0
	}
	last := end
	if r.HeldAt != nil && last.After(*r.HeldAt) {
		last = *r.HeldAt
	}
	return min(hours(r.used(end).Hours()*r.Rate*t.factor(*r.StartedAt, last)), r.reserve())
}

func (r Run) input() RunInput {
	in := RunInput{
		Kind: r.Kind, Title: r.Title, Command: r.Command, TemplateID: r.TemplateID, ClassID: r.ClassID, MaxHours: r.MaxHours,
	}
	if r.WorkKind != nil && r.WorkID != nil {
		in.WorkKind, in.WorkID = *r.WorkKind, *r.WorkID
	}
	if r.TeamID != nil {
		in.TeamID = *r.TeamID
	}
	return in
}

func (r Run) interactive() bool {
	return r.Kind == RunJupyter || r.Kind == RunVSCode
}

func (r Run) seenBy(userID string) Run {
	if r.UserID != userID {
		r.URL = ""
	}
	return r
}

func (r Run) workspace() string {
	kind, id := r.owner()
	return kind + "-" + id
}
