package models

import "time"

type Lesson struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ModulesID int64     `json:"module_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Units     []Unit    `gorm:"foreignKey:LessonsID"`
}
