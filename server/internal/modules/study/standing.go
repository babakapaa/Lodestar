package study

type Standing struct {
	Role   Role
	Manage bool
}

func (s Standing) Visible() bool {
	return s.Role != RoleNone || s.Manage
}

func (s Standing) Teaches() bool {
	return s.Role == RoleTeacher || s.Manage
}

func (s Standing) Shown() Role {
	if s.Role == RoleNone && s.Manage {
		return RoleManager
	}
	return s.Role
}
