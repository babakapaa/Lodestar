package access

import "time"

type Assignment struct {
	ID        string     `json:"id"`
	RoleID    string     `json:"roleId"`
	UserID    string     `json:"userId"`
	Scope     string     `json:"scope"`
	ExpiresAt *time.Time `json:"expiresAt"`
	GrantedBy *string    `json:"grantedBy"`
	GrantedAt time.Time  `json:"grantedAt"`
}
