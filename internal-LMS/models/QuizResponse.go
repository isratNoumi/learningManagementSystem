package models

import "time"

type QuizTaker struct {
	UserID       int64          `json:"user_id"`
	CoursesID    int64          `json:"courses_id"`
	QuizResponse []QuizResponse `json:"answers"`
}
type QuizResponse struct {
	UnitsDetailsID int64  `json:"units_details_id"`
	Response       string `json:"response"`
}
type Response struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	UsersID        int64     `json:"user_id"`
	CoursesID      int64     `json:"courses_id"`
	UnitsDetailsID int64     `json:"units_details_id"`
	Response       string    `json:"response"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
type AnswerFeedback struct {
	UnitsDetailsID int64  `json:"units_details_id"`
	UserAnswer     string `json:"user_answer"`
	CorrectAnswer  string `json:"correct_answer"`
	IsCorrect      bool   `json:"is_correct"`
}
