package attendance

import "errors"

var ErrLessonNotFound = errors.New("lesson not found")

var ErrLessonCancelled = errors.New("lesson is cancelled")

var ErrSessionNotFound = errors.New("attendance session not found")

var ErrSessionClosed = errors.New("attendance session is closed")

var ErrInvalidMark = errors.New("attendance mark is invalid")

var ErrNotInLesson = errors.New("not a student of this lesson")

var ErrInvalidCheck = errors.New("check-in is invalid")

var ErrMethodOff = errors.New("this check-in method is turned off")

var ErrOffCampus = errors.New("check-in works only from the university network")

var ErrWrongCode = errors.New("code or qr is wrong or expired")

var ErrTooManyTries = errors.New("too many wrong codes")

var ErrInvalidRange = errors.New("date range is invalid")
