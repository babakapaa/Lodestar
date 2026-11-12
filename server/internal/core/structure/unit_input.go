package structure

type UnitInput struct {
	ParentID *string `json:"parentId"`
	Kind     string  `json:"kind"`
	Name     string  `json:"name"`
	Code     *string `json:"code"`
}
