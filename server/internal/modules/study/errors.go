package study

import "errors"

var ErrCourseNotFound = errors.New("course not found")

var ErrInvalidCourse = errors.New("course title, code or description is invalid")

var ErrCodeTaken = errors.New("course code already taken")

var ErrArchived = errors.New("course is archived")

var ErrUnitNotFound = errors.New("unit not found")

var ErrUserNotFound = errors.New("user not found")
