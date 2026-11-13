package module

type PilotCourse struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Code  *string `json:"code"`
}
