package attendance

type Counts struct {
	Present  int `json:"present"`
	Late     int `json:"late"`
	Absent   int `json:"absent"`
	Excused  int `json:"excused"`
	Unmarked int `json:"unmarked"`
}

func (c *Counts) add(status *MarkStatus) {
	if status == nil {
		c.Unmarked++
		return
	}
	switch *status {
	case Present:
		c.Present++
	case Late:
		c.Late++
	case Absent:
		c.Absent++
	case Excused:
		c.Excused++
	}
}
