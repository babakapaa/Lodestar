package profile

type Profile struct {
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName"`
	About       string    `json:"about"`
	Contacts    []Contact `json:"contacts"`
}
