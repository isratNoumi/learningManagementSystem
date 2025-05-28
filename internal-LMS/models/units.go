package models

import "time"

type Unit struct {
	ID          int64         `json:"id,omitempty"`
	LessonsID   int64         `json:"lessons_id,omitempty"`
	Type        string        `json:"type,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	UnitsFields []UnitsFields `json:"units_fields,omitempty" gorm:"foreignKey:UnitsID"`
}

type UnitsFields struct {
	ID           int64          `json:"id,omitempty"`
	UnitsID      int64          `json:"units_id,omitempty"`
	Fields       string         `json:"fields,omitempty"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	UnitsDetails []UnitsDetails `json:"units_details,omitempty" gorm:"foreignKey:UnitsFieldsID"`
}

type UnitsDetails struct {
	ID            int64         `json:"id,omitempty"`
	UnitsFieldsID int64         `json:"units_fields_id,omitempty"`
	Description   string        `json:"description,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
	QuizType      []QuizTypes   `json:"quiz_types,omitempty" gorm:"foreignKey:UnitsDetailsID"`
	QuizOption    []QuizOptions `json:"quiz_option,omitempty" gorm:"foreignKey:UnitsDetailsID"`
}
type Video struct {
	URL      string `json:"url,omitempty"`
	Duration string `json:"duration,omitempty"`
}
type Content struct {
	Text   string `json:"text,omitempty"`
	Length string `json:"length,omitempty"`
}
type UnitsViews struct {
	ID       int64     `json:"id,omitempty"`
	UsersID  int64     `json:"user_id,omitempty"`
	UnitsID  int64     `json:"units_id,omitempty"`
	ViewedAt time.Time `json:"viewed_at" gorm:"autoCreateTime; <-:create"` // Automatically set the time when the view is created

}
