package access

type RoleInput struct {
	Name        string       `json:"name"`
	Color       string       `json:"color"`
	Icon        string       `json:"icon"`
	Permissions []Permission `json:"permissions"`
}
