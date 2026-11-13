package chat

const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

func manages(role string) bool {
	return role == RoleOwner || role == RoleAdmin
}
