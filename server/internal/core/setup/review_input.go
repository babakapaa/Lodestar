package setup

type ReviewInput struct {
	Approve bool   `json:"approve"`
	Note    string `json:"note"`
}
