package attendance

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const (
	CodeLessonNotFound  = "lesson_not_found"
	CodeLessonCancelled = "lesson_cancelled"
	CodeSessionNotFound = "attendance_session_not_found"
	CodeSessionClosed   = "attendance_session_closed"
	CodeInvalidMark     = "invalid_attendance_mark"
	CodeNotInLesson     = "not_in_lesson"
	CodeInvalidCheck    = "invalid_check_in"
	CodeMethodOff       = "check_in_method_off"
	CodeOffCampus       = "off_campus"
	CodeWrongCode       = "wrong_attendance_code"
	CodeTooManyTries    = "too_many_attendance_tries"
	CodeInvalidRange    = "invalid_range"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrLessonNotFound, http.StatusNotFound, CodeLessonNotFound},
	{ErrLessonCancelled, http.StatusConflict, CodeLessonCancelled},
	{ErrSessionNotFound, http.StatusNotFound, CodeSessionNotFound},
	{ErrSessionClosed, http.StatusConflict, CodeSessionClosed},
	{ErrInvalidMark, http.StatusBadRequest, CodeInvalidMark},
	{ErrNotInLesson, http.StatusForbidden, CodeNotInLesson},
	{ErrInvalidCheck, http.StatusBadRequest, CodeInvalidCheck},
	{ErrMethodOff, http.StatusForbidden, CodeMethodOff},
	{ErrOffCampus, http.StatusForbidden, CodeOffCampus},
	{ErrWrongCode, http.StatusBadRequest, CodeWrongCode},
	{ErrTooManyTries, http.StatusTooManyRequests, CodeTooManyTries},
	{ErrInvalidRange, http.StatusBadRequest, CodeInvalidRange},
}

func writeError(w http.ResponseWriter, err error) {
	for _, e := range errorCodes {
		if errors.Is(err, e.err) {
			apierror.Write(w, e.status, e.code)
			return
		}
	}
	access.WriteError(w, err)
}

func badRequest(w http.ResponseWriter) {
	apierror.Write(w, http.StatusBadRequest, access.CodeInvalidRequest)
}

func respond(w http.ResponseWriter, body any, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, body)
}
