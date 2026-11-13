package study

import "time"

type UpcomingDefense struct {
	AssignmentID string
	StudentID    string
	SlotID       string
	Assignment   string
	Course       string
	StartsAt     time.Time
}
