package lti

type SourceInput struct {
	Source *string `json:"source"`
}

func (in SourceInput) valid() bool {
	return in.Source == nil || *in.Source == SourceLodestar || *in.Source == SourceLMS
}
