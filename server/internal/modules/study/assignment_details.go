package study

import "github.com/babakapaa/lodestar/server/internal/core/files"

type AssignmentDetails struct {
	Assignment
	CourseTitle string       `json:"courseTitle"`
	Role        Role         `json:"role"`
	Files       []files.File `json:"files"`
	Mine        *StudentWork `json:"mine"`
	Counts      *Counts      `json:"counts"`
}
