package models

import "time"

type Lesson struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ModulesID int64     `json:"module_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Units     []Unit    `gorm:"foreignKey:LessonsID"`
}
