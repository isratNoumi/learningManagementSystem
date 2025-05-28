package models

import "time"

type MCQuiz struct {
	LessonName   string `json:"name"`
	UnitType     string `json:"type"`
	QuestionType string `json:"question_type"`
	ContentType  string `json:"quiz_type"`
	Question     string `json:"question"`
}
type QuizTypes struct {
	ID             int64     `gorm:"primaryKey"`
	UnitsDetailsID int64     `json:"units_details_id"`
	Type           string    `json:"type"`
	CreatedAt      time.Time `json:"created_at ,omitempty" gorm:"autoCreateTime<-:create"`
	UpdatedAt      time.Time `json:"updated_at ,omitempty" gorm:"autoUpdateTime"`
}

type QuizOptions struct {
	ID             int64     `gorm:"primaryKey"`
	UnitsDetailsID int64     `json:"units_details_id"`
	Options        string    `json:"options"` // Assuming options is a string column
	CreatedAt      time.Time `json:"created_at ,omitempty" gorm:"autoCreateTime<-:create"`
	UpdatedAt      time.Time `json:"updated_at ,omitempty" gorm:"autoUpdateTime"`
}
