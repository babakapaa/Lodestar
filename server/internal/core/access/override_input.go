package access

type OverrideInput struct {
	Allow []Permission `json:"allow"`
	Deny  []Permission `json:"deny"`
}
