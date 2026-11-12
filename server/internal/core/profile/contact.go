package profile

type Contact struct {
	Kind   string `json:"kind"`
	Value  string `json:"value"`
	Public bool   `json:"public"`
}
