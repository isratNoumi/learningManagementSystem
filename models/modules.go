package models

import "time"

type Module struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CoursesID int64     `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Lessons   []Lesson  `gorm:"foreignKey:ModulesID"`
}
