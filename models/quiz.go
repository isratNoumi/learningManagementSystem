package models

type MCQuiz struct {
	LessonName   string `json:"name"`
	UnitType     string `json:"type"`
	QuestionType string `json:"question_type"`
	ContentType  string `json:"quiz_type"`
	Question     string `json:"question"`
}
type QuizTypes struct {
	ID             int64  `gorm:"primaryKey"`
	UnitsDetailsID int64  `json:"units_details_id"`
	Type           string `json:"type"`
}

type QuizOptions struct {
	ID             uint   `gorm:"primaryKey"`
	UnitsDetailsID uint   `json:"units_details_id"`
	Options        string `json:"options"` // Assuming options is a string column
}
