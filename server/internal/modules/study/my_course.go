package study

type MyCourse struct {
	Course
	Role Role `json:"role"`
}
