package module

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrLessonNotFound = errors.New("lesson not found")

type LessonSource interface {
	Lesson(ctx context.Context, tx pgx.Tx, seriesID, date string) (LessonRef, error)
	Course(ctx context.Context, tx pgx.Tx, courseID, from, to string) (LessonCourse, error)
}
