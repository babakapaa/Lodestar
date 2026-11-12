package study

type AssignmentItem struct {
	Assignment
	Mine   *Progress `json:"mine"`
	Counts *Counts   `json:"counts"`
}
