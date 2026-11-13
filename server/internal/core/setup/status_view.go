package setup

type StatusView struct {
	Launched bool     `json:"launched"`
	Manages  bool     `json:"manages"`
	Assigned []string `json:"assigned"`
	Pilot    bool     `json:"pilot"`
	Included bool     `json:"included"`
}
