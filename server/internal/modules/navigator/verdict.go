package navigator

type Verdict struct {
	Text      string
	Citations []Citation
	Leak      bool
	Removed   int
	Articles  []Article
}

func (v Verdict) Passed() bool {
	return v.Text != "" && !v.Leak
}
