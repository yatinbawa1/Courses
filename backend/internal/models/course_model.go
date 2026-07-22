package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	CourseID          uuid.UUID `json:"course_id"`
	CourseName        string    `json:"course_name"`
	CourseDescription string    `json:"course_description"`
	CreationDate      time.Time `json:"creation_date"`
	CourseThumbnail   string 	`json:"course_thumbnail"`
}

type AssetType string

const (
	AssetVideo AssetType = "VIDEO"
	AssetPDF   AssetType = "PDF"
	AssetAudio AssetType = "AUDIO"
	AssetImage AssetType = "IMAGE"
	AssetQuiz  AssetType = "QUIZ"
)

type CourseChapter struct {
	ChapterContentType AssetType
	// Todo
}
