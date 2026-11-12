package notes

type Role string

const (
	RoleOwner  Role = "owner"
	RoleEditor Role = "editor"
	RoleReader Role = "reader"
)

func (r Role) canEdit() bool {
	return r == RoleOwner || r == RoleEditor
}
