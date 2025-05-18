package models

import "time"

type Course struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	TotalScore float64   `json:"total_score"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Modules    []Module  `gorm:"foreignKey:CoursesID"`
	Links      []Link    `json:"links" gorm:"-"`
}

type CourseDTO struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	TotalScore float64   `json:"total_score"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
