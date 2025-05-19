package models

import "time"

type Module struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CoursesID int64     `json:"course_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Lessons   []Lesson  `gorm:"foreignKey:ModulesID"`
}
