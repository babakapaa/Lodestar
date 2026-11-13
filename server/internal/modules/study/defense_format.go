package study

type DefenseFormat string

const (
	DefenseFull   DefenseFormat = "full"
	DefenseShort  DefenseFormat = "short"
	DefenseSample DefenseFormat = "sample"
)

func (f DefenseFormat) valid() bool {
	return f == DefenseFull || f == DefenseShort || f == DefenseSample
}
