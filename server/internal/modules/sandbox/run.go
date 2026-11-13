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

func (r Run) cost(end time.Time) float64 {
	if r.StartedAt == nil || !end.After(*r.StartedAt) {
		return 0
	}
	return min(hours(end.Sub(*r.StartedAt).Hours()*r.Rate), r.reserve())
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
