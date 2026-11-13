package stats

type Report struct {
	ParentID  *string     `json:"parentId"`
	Days      int         `json:"days"`
	Threshold int         `json:"threshold"`
	Full      bool        `json:"full"`
	Units     []UnitStats `json:"units"`
}
