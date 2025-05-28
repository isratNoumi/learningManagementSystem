package models

import "time"

type Lesson struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ModulesID int64     `json:"module_id"`
	CreatedAt time.Time `json:"created-at,omitempty" gorm:"autoCreateTime<-:create"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	Units     []Unit    `json:"units" gorm:"foreignKey:LessonsID"`
}
