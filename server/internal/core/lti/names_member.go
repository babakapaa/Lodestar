package lti

type NamesMember struct {
	UserID     string   `json:"user_id"`
	Status     string   `json:"status"`
	Roles      []string `json:"roles"`
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	GivenName  string   `json:"given_name"`
	FamilyName string   `json:"family_name"`
}

func (m NamesMember) Active() bool {
	return m.Status == "" || m.Status == "Active"
}

func (m NamesMember) Person() Person {
	c := LaunchClaims{Name: m.Name, GivenName: m.GivenName, FamilyName: m.FamilyName}
	return Person{Subject: m.UserID, Email: m.Email, Name: c.DisplayName()}
}

func (m NamesMember) Instructor() bool {
	return LaunchClaims{Roles: m.Roles}.Instructor()
}
