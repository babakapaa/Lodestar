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

var ErrInvalidSearchMap = errors.New("search map is malformed")

var ErrNoDefense = errors.New("assignment has no defense")

var ErrInvalidDefense = errors.New("defense questions or protocol are invalid")

var ErrInvalidSlot = errors.New("defense slot is invalid")

var ErrSlotNotFound = errors.New("defense slot not found")

var ErrSlotFull = errors.New("defense slot is full")

var ErrSlotStarted = errors.New("defense slot has already started")

var ErrDefenseFinished = errors.New("defense is already finished")

var ErrNotChosen = errors.New("student is not chosen for this defense")

var ErrNoWork = errors.New("student has not submitted work yet")

var ErrAIUnavailable = errors.New("language model is not configured")

var ErrAIFailed = errors.New("language model did not give a usable answer")

var ErrSearchMapIncomplete = errors.New("search map does not meet the assignment rules")

var ErrTeamWork = errors.New("this work is handed in by the team")

var ErrBadImport = errors.New("import file is malformed")
