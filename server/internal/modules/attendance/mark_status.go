package attendance

type MarkStatus string

const (
	Present MarkStatus = "present"
	Late    MarkStatus = "late"
	Absent  MarkStatus = "absent"
	Excused MarkStatus = "excused"
)

func (s MarkStatus) valid() bool {
	switch s {
	case Present, Late, Absent, Excused:
		return true
	}
	return false
}
