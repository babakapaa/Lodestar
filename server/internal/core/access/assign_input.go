package access

import "time"

type AssignInput struct {
	RoleID    string     `json:"roleId"`
	UserID    string     `json:"userId"`
	Scope     string     `json:"scope"`
	ExpiresAt *time.Time `json:"expiresAt"`
}
