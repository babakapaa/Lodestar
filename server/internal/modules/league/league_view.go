package league

type LeagueView struct {
	ID          string          `json:"id"`
	Kind        string          `json:"kind"`
	OriginID    string          `json:"originId"`
	Title       string          `json:"title"`
	Role        Role            `json:"role"`
	Season      Season          `json:"season"`
	Divisions   Divisions       `json:"divisions"`
	Confirm     bool            `json:"confirm"`
	Formula     []WorkWeight    `json:"formula"`
	Me          *Standing       `json:"me"`
	Table       []TableRow      `json:"table"`
	Entries     []Entry         `json:"entries"`
	Transitions []Transition    `json:"transitions"`
	Disputes    []Dispute       `json:"disputes"`
	Players     []PlayerRow     `json:"players"`
	Pending     []TransitionRow `json:"pending"`
	Open        []DisputeRow    `json:"openDisputes"`
}
