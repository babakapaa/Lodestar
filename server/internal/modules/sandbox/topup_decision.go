package sandbox

type TopupDecision struct {
	Approve bool     `json:"approve"`
	Hours   *float64 `json:"hours"`
}

func (d TopupDecision) granted(t Topup) (float64, error) {
	if !d.Approve {
		return 0, nil
	}
	if d.Hours == nil {
		return t.Hours, nil
	}
	h := hours(*d.Hours)
	if h <= 0 || h > 10000 {
		return 0, ErrInvalidTopup
	}
	return h, nil
}
