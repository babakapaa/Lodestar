package sandbox

type EstimateInput struct {
	Kind       string  `json:"kind"`
	TemplateID string  `json:"templateId"`
	MaxHours   float64 `json:"maxHours"`
	TeamID     string  `json:"teamId"`
}

func (in EstimateInput) normalized() (EstimateInput, error) {
	in.MaxHours = defaultHours(in.Kind, hours(in.MaxHours))
	if in.Kind != RunJupyter && in.Kind != RunVSCode && in.Kind != RunJob && in.Kind != RunDemo ||
		in.TemplateID == "" || in.MaxHours <= 0 || in.MaxHours > 168 {
		return in, ErrInvalidRun
	}
	return in, nil
}
