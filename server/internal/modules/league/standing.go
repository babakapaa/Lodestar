package league

type Standing struct {
	Points   int     `json:"points"`
	Division int     `json:"division"`
	ToNext   *int    `json:"toNext"`
	Growth   int     `json:"growth"`
	Display  Display `json:"display"`
	Nickname string  `json:"nickname"`
}
