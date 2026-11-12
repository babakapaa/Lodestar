package tenant

import "time"

type Tenant struct {
	ID        string
	Slug      string
	Name      string
	CreatedAt time.Time
}
