package team

type LeaveResult struct {
	Left    bool          `json:"left"`
	Request *LeaveRequest `json:"request"`
}
