package attendance

type StudentRecord struct {
	UserID     string   `json:"userId"`
	Name       string   `json:"name"`
	Counts     Counts   `json:"counts"`
	Rate       *float64 `json:"rate"`
	Streak     int      `json:"streak"`
	LongAbsent bool     `json:"longAbsent"`
}
