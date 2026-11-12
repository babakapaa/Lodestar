package demo

import "github.com/babakapaa/lodestar/server/internal/platform/demoidp"

var (
	Student = demoidp.User{Subject: "demo-student", Email: "student@demo.invalid", Name: "Айгерим Сапарова", Title: "Студентка"}
	Teacher = demoidp.User{Subject: "demo-teacher", Email: "teacher@demo.invalid", Name: "Ерлан Жумабаев", Title: "Преподаватель"}
	Admin   = demoidp.User{Subject: "demo-admin", Email: "admin@demo.invalid", Name: "Дана Ахметова", Title: "Администратор вуза"}
)

func People() []demoidp.User {
	return []demoidp.User{Student, Teacher, Admin}
}
