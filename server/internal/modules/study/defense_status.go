package study

type DefenseStatus string

const (
	DefensePlanned DefenseStatus = "planned"
	DefenseHeld    DefenseStatus = "held"
	DefenseMissed  DefenseStatus = "missed"
)

func (s DefenseStatus) finished() bool {
	return s == DefenseHeld || s == DefenseMissed
}
