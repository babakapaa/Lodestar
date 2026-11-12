package structure

type UnitDetails struct {
	Unit
	Path []Unit `json:"path"`
}
