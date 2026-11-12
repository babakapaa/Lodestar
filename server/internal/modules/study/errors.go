package study

import "errors"

var ErrCourseNotFound = errors.New("course not found")

var ErrInvalidCourse = errors.New("course title, code or description is invalid")

var ErrCodeTaken = errors.New("course code already taken")

var ErrArchived = errors.New("course is archived")

var ErrUnitNotFound = errors.New("unit not found")

var ErrUserNotFound = errors.New("user not found")

var ErrAssignmentNotFound = errors.New("assignment not found")

var ErrInvalidAssignment = errors.New("assignment title, grading or files are invalid")

var ErrHasWork = errors.New("assignment already has submitted work")

var ErrNotPublished = errors.New("assignment is not published")

var ErrInvalidSubmission = errors.New("submission is empty or malformed")

var ErrInvalidGrade = errors.New("grade does not match the assignment grading")

var ErrNotStudent = errors.New("user is not a student of the course")

var ErrBadFile = errors.New("file is not ready or belongs to someone else")

var ErrSeriesNotFound = errors.New("schedule series not found")

var ErrInvalidSeries = errors.New("schedule series is invalid")

var ErrInvalidChange = errors.New("lesson change is invalid")

var ErrInvalidRange = errors.New("date range is invalid")
