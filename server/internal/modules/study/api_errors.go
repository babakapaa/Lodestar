package study

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

const (
	CodeCourseNotFound = "course_not_found"
	CodeInvalidCourse  = "invalid_course"
	CodeCodeTaken      = "course_code_taken"
	CodeArchived       = "course_archived"
	CodeUnitNotFound   = "unit_not_found"
	CodeUserNotFound   = "user_not_found"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrCourseNotFound, http.StatusNotFound, CodeCourseNotFound},
	{ErrInvalidCourse, http.StatusBadRequest, CodeInvalidCourse},
	{ErrCodeTaken, http.StatusConflict, CodeCodeTaken},
	{ErrArchived, http.StatusConflict, CodeArchived},
	{ErrUnitNotFound, http.StatusNotFound, CodeUnitNotFound},
	{ErrUserNotFound, http.StatusNotFound, CodeUserNotFound},
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
