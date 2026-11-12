package access

type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Color       string       `json:"color"`
	Icon        string       `json:"icon"`
	Position    int          `json:"position"`
	Permissions []Permission `json:"permissions"`
	Default     bool         `json:"default"`
}
