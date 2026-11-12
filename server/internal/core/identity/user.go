package identity

import "time"

type User struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenantId"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	CreatedAt   time.Time `json:"createdAt"`
}
