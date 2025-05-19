package models

import "time"

type Unit struct {
	ID           int64          `json:"id,omitempty"`
	LessonsID    int64          `json:"lesson_id,omitempty"`
	Type         string         `json:"type,omitempty"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	UnitsFields  []UnitsFields  `json:"units_fields" gorm:"foreignKey:UnitsID"`
	UnitsDetails []UnitsDetails `json:"-" gorm:"foreignKey:UnitsID"`
}
type UnitsFields struct {
	ID           int64          `json:"id,omitempty"`
	UnitsID      int64          `json:"units_id,omitempty"`
	UnitsFields  string         `json:"unit_fields,omitempty"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	UnitsDetails []UnitsDetails `json:"units_details" gorm:"foreignKey:UnitsFieldsID"`
}

type UnitsDetails struct {
	ID            int64         `json:"id,omitempty"`
	UnitsID       int64         `json:"units_id,omitempty"`
	UnitsFieldsID int64         `json:"unit_fields,omitempty"`
	Description   string        `json:"description,omitempty"`
	CreatedAt     time.Time     `json:"-,omitempty"`
	UpdatedAt     time.Time     `json:"-,omitempty"`
	QuizType      []QuizTypes   `json:"quiz_types,omitempty" gorm:"foreignKey:UnitsDetailsID"`
	QuizOption    []QuizOptions `json:"quiz_option,omitempty" gorm:"foreignKey:UnitsDetailsID"`
}
type Video struct {
	URL      string `json:"url"`
	Duration int    `json:"duration"`
}
type MCQQuiz struct {
	Question   string       `json:"question"`
	Answer     string       `json:"answer"`
	MCQOptions []MCQOptions `json:"mcq_options"`
}
type MCQOptions struct {
	OptionA string `json:"optionA"`
	OptionB string `json:"optionB"`
	OptionC string `json:"optionC"`
	OptionD string `json:"optionD"`
}
type TrueFalseOption struct {
	OptionA string `json:"optionA"`
	OptionB string `json:"optionB"`
}
