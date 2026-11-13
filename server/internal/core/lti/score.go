package lti

import "time"

type Score struct {
	UserID           string    `json:"userId"`
	ScoreGiven       float64   `json:"scoreGiven"`
	ScoreMaximum     float64   `json:"scoreMaximum"`
	ActivityProgress string    `json:"activityProgress"`
	GradingProgress  string    `json:"gradingProgress"`
	Timestamp        time.Time `json:"timestamp"`
}
