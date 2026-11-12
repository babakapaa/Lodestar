package schema

type Type string

const (
	Bool     Type = "bool"
	Integer  Type = "integer"
	Text     Type = "text"
	Choice   Type = "choice"
	TimeZone Type = "time_zone"
)
