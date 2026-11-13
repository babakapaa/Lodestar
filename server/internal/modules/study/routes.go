package study

import "net/http"

type Routes struct {
	Service     Service
	Assignments AssignmentService
	Schedule    ScheduleService
	Defenses    DefenseService
	Gradebook   GradebookService
	Imports     StudentImportService
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Service
	mux.Handle("GET /v1/study/courses", MyCoursesHandler{service: s})
	mux.Handle("POST /v1/study/courses", CreateCourseHandler{service: s})
	mux.Handle("GET /v1/study/catalog", CatalogHandler{service: s})
	mux.Handle("GET /v1/study/courses/{id}", CourseHandler{service: s})
	mux.Handle("PUT /v1/study/courses/{id}", UpdateCourseHandler{service: s})
	mux.Handle("POST /v1/study/courses/{id}/archive", ArchiveCourseHandler{service: s, archived: true})
	mux.Handle("POST /v1/study/courses/{id}/restore", ArchiveCourseHandler{service: s})
	mux.Handle("GET /v1/study/courses/{id}/students", StudentsHandler{service: s})
	mux.Handle("POST /v1/study/courses/{id}/students/import", StudentImportHandler{service: r.Imports})
	mux.Handle("GET /v1/study/courses/{id}/gradebook", GradebookHandler{service: r.Gradebook})
	a := r.Assignments
	mux.Handle("GET /v1/study/courses/{id}/assignments", AssignmentsHandler{service: a})
	mux.Handle("POST /v1/study/courses/{id}/assignments", CreateAssignmentHandler{service: a})
	mux.Handle("GET /v1/study/assignments/{id}", AssignmentHandler{service: a})
	mux.Handle("PUT /v1/study/assignments/{id}", UpdateAssignmentHandler{service: a})
	mux.Handle("DELETE /v1/study/assignments/{id}", DeleteAssignmentHandler{service: a})
	mux.Handle("POST /v1/study/assignments/{id}/submissions", SubmitHandler{service: a})
	mux.Handle("GET /v1/study/assignments/{id}/work", WorkHandler{service: a})
	mux.Handle("GET /v1/study/assignments/{id}/search-map", SearchMapHandler{service: a})
	mux.Handle("PUT /v1/study/assignments/{id}/search-map", SaveSearchMapHandler{service: a})
	mux.Handle("GET /v1/study/assignments/{id}/students/{student}/submissions", AttemptsHandler{service: a})
	mux.Handle("PUT /v1/study/assignments/{id}/grades/{student}", GradeHandler{service: a})
	sc := r.Schedule
	mux.Handle("GET /v1/study/courses/{id}/series", SeriesListHandler{service: sc})
	mux.Handle("POST /v1/study/courses/{id}/series", CreateSeriesHandler{service: sc})
	mux.Handle("PUT /v1/study/series/{id}", UpdateSeriesHandler{service: sc})
	mux.Handle("DELETE /v1/study/series/{id}", DeleteSeriesHandler{service: sc})
	mux.Handle("PUT /v1/study/series/{id}/changes/{date}", ChangeHandler{service: sc})
	mux.Handle("DELETE /v1/study/series/{id}/changes/{date}", RevertChangeHandler{service: sc})
	mux.Handle("GET /v1/study/schedule", ScheduleHandler{service: sc})
	mux.Handle("GET /v1/study/courses/{id}/schedule", CourseScheduleHandler{service: sc})
	mux.Handle("GET /v1/study/agenda", AgendaHandler{service: sc})
	r.mountDefenses(mux)
	for _, list := range []List{ListTeachers, ListUnits, ListStudents} {
		path := "/v1/study/courses/{id}/" + string(list) + "/{member}"
		mux.Handle("PUT "+path, RosterChangeHandler{service: s, list: list, add: true})
		mux.Handle("DELETE "+path, RosterChangeHandler{service: s, list: list})
	}
}
