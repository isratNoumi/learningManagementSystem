package models

import "time"

type Unit struct {
	ID          int64         `json:"id,omitempty"`
	LessonsID   int64         `json:"lesson_id,omitempty"`
	Type        string        `json:"type,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	UnitsFields []UnitsFields `json:"unit_fields" gorm:"foreignKey:UnitsID"`
}

type UnitsFields struct {
	ID         int64     `json:"id,omitempty"`
	UnitsID    int64     `json:"unit_id,omitempty"`
	UnitFields string    `json:"unit_fields,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type UnitsDetails struct {
	ID            int64     `json:"id,omitempty"`
	UnitsID       int64     `json:"units_id,omitempty"`
	UnitsFieldsID int64     `json:"unit_fields,omitempty"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
