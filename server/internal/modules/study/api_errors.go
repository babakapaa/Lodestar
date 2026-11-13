package study

import (
	"errors"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	CodeCourseNotFound     = "course_not_found"
	CodeInvalidCourse      = "invalid_course"
	CodeCodeTaken          = "course_code_taken"
	CodeArchived           = "course_archived"
	CodeUnitNotFound       = "unit_not_found"
	CodeUserNotFound       = "user_not_found"
	CodeAssignmentNotFound = "assignment_not_found"
	CodeInvalidAssignment  = "invalid_assignment"
	CodeHasWork            = "assignment_has_work"
	CodeNotPublished       = "assignment_not_published"
	CodeInvalidSubmission  = "invalid_submission"
	CodeInvalidGrade       = "invalid_grade"
	CodeNotStudent         = "not_course_student"
	CodeBadFile            = "file_not_usable"
	CodeSeriesNotFound     = "series_not_found"
	CodeInvalidSeries      = "invalid_series"
	CodeInvalidChange      = "invalid_change"
	CodeInvalidRange       = "invalid_range"
	CodeInvalidSearchMap   = "invalid_search_map"
	CodeMapIncomplete      = "search_map_incomplete"
	CodeNoDefense          = "assignment_has_no_defense"
	CodeInvalidDefense     = "invalid_defense"
	CodeInvalidSlot        = "invalid_defense_slot"
	CodeSlotNotFound       = "defense_slot_not_found"
	CodeSlotFull           = "defense_slot_full"
	CodeSlotStarted        = "defense_slot_started"
	CodeDefenseFinished    = "defense_finished"
	CodeNotChosen          = "defense_not_chosen"
	CodeNoWork             = "no_submitted_work"
	CodeTeamWork           = "team_work"
	CodeAIUnavailable      = "study_ai_unavailable"
	CodeAIFailed           = "study_ai_failed"
	CodeBadImport          = "bad_import"
)

var errorCodes = []struct {
	err    error
	status int
	code   string
}{
	{ErrBadImport, http.StatusBadRequest, CodeBadImport},
	{ErrCourseNotFound, http.StatusNotFound, CodeCourseNotFound},
	{ErrInvalidCourse, http.StatusBadRequest, CodeInvalidCourse},
	{ErrCodeTaken, http.StatusConflict, CodeCodeTaken},
	{ErrArchived, http.StatusConflict, CodeArchived},
	{ErrUnitNotFound, http.StatusNotFound, CodeUnitNotFound},
	{ErrUserNotFound, http.StatusNotFound, CodeUserNotFound},
	{ErrAssignmentNotFound, http.StatusNotFound, CodeAssignmentNotFound},
	{ErrInvalidAssignment, http.StatusBadRequest, CodeInvalidAssignment},
	{ErrHasWork, http.StatusConflict, CodeHasWork},
	{ErrNotPublished, http.StatusConflict, CodeNotPublished},
	{ErrInvalidSubmission, http.StatusBadRequest, CodeInvalidSubmission},
	{ErrInvalidGrade, http.StatusBadRequest, CodeInvalidGrade},
	{ErrNotStudent, http.StatusNotFound, CodeNotStudent},
	{ErrBadFile, http.StatusBadRequest, CodeBadFile},
	{ErrSeriesNotFound, http.StatusNotFound, CodeSeriesNotFound},
	{ErrInvalidSeries, http.StatusBadRequest, CodeInvalidSeries},
	{ErrInvalidChange, http.StatusBadRequest, CodeInvalidChange},
	{ErrInvalidRange, http.StatusBadRequest, CodeInvalidRange},
	{ErrInvalidSearchMap, http.StatusBadRequest, CodeInvalidSearchMap},
	{ErrSearchMapIncomplete, http.StatusUnprocessableEntity, CodeMapIncomplete},
	{ErrNoDefense, http.StatusConflict, CodeNoDefense},
	{ErrInvalidDefense, http.StatusBadRequest, CodeInvalidDefense},
	{ErrInvalidSlot, http.StatusBadRequest, CodeInvalidSlot},
	{ErrSlotNotFound, http.StatusNotFound, CodeSlotNotFound},
	{ErrSlotFull, http.StatusConflict, CodeSlotFull},
	{ErrSlotStarted, http.StatusConflict, CodeSlotStarted},
	{ErrDefenseFinished, http.StatusConflict, CodeDefenseFinished},
	{ErrNotChosen, http.StatusForbidden, CodeNotChosen},
	{ErrNoWork, http.StatusConflict, CodeNoWork},
	{ErrTeamWork, http.StatusConflict, CodeTeamWork},
	{ErrAIUnavailable, http.StatusServiceUnavailable, CodeAIUnavailable},
	{ErrAIFailed, http.StatusBadGateway, CodeAIFailed},
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

func problem(err error) error {
	for _, e := range errorCodes {
		if errors.Is(err, e.err) {
			return module.TeamWorkProblem{Status: e.status, Code: e.code}
		}
	}
	return err
}
