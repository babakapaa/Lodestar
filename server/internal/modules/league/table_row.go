package league

type TableRow struct {
	Place     int     `json:"place"`
	Name      string  `json:"name"`
	Display   Display `json:"display"`
	Points    int     `json:"points"`
	Growth    int     `json:"growth"`
	Me        bool    `json:"me"`
	TopGrowth bool    `json:"topGrowth"`
}
