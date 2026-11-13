package sandbox

type TransferInput struct {
	TeamID string  `json:"teamId"`
	Hours  float64 `json:"hours"`
}

func (in TransferInput) normalized() (TransferInput, error) {
	in.Hours = hours(in.Hours)
	if in.TeamID == "" || in.Hours <= 0 || in.Hours > 10000 {
		return in, ErrInvalidTransfer
	}
	return in, nil
}
