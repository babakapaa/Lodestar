package league

import "time"

type Grant struct {
	ID        string    `json:"id"`
	AwardID   string    `json:"awardId"`
	Title     string    `json:"title"`
	GPUHours  int       `json:"gpuHours"`
	UserID    string    `json:"userId"`
	SeasonID  string    `json:"seasonId"`
	Reason    string    `json:"reason"`
	GrantedAt time.Time `json:"grantedAt"`
}
