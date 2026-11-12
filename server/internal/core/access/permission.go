package access

type Permission string

const (
	Administrator Permission = "core.administrator"
	ManageRoles   Permission = "core.roles.manage"
)

func corePermissions() []Permission {
	return []Permission{Administrator, ManageRoles}
}
