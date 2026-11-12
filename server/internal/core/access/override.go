package access

type Override struct {
	RoleID string       `json:"roleId"`
	Allow  []Permission `json:"allow"`
	Deny   []Permission `json:"deny"`
}
