package attendance

import "net/http"

type Routes struct {
	Service Service
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Service
	mux.Handle("POST /v1/attendance/lessons/{series}/{date}/open", OpenHandler{service: s})
	mux.Handle("GET /v1/attendance/sessions/{id}", SessionHandler{service: s})
	mux.Handle("GET /v1/attendance/sessions/{id}/qr", QRHandler{service: s})
	mux.Handle("POST /v1/attendance/sessions/{id}/code", CodeHandler{service: s})
	mux.Handle("PUT /v1/attendance/sessions/{id}/marks/{user}", MarkHandler{service: s})
	mux.Handle("POST /v1/attendance/sessions/{id}/close", CloseHandler{service: s})
	mux.Handle("POST /v1/attendance/check", CheckHandler{service: s})
	mux.Handle("GET /v1/attendance/open", OpenLessonsHandler{service: s})
	mux.Handle("GET /v1/attendance/mine", MineHandler{service: s})
	mux.Handle("GET /v1/attendance/courses/{id}", CourseHandler{service: s})
	mux.Handle("GET /v1/attendance/courses/{id}/export", ExportHandler{service: s})
}
