package models

import "time"

type Course struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	TotalScore float64   `json:"total_score"`
	CreatedAt  time.Time `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	Modules    []Module  `json:"modules,omitempty" gorm:"foreignKey:CoursesID"`
	Links      []Link    `json:"links,omitempty" gorm:"-"`
}

type CourseDTO struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	TotalScore float64   `json:"total_score"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type ModuleResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CoursesID int    `json:"courses_id"`
}
type InstructorResponse struct {
	CoursesName     string `json:"courses_name"`
	NoOfStudents    int    `json:"no_of_students"`
	CoursesCategory string `json:"courses_category"`
}
type UserResponse struct {
	CoursesName    string `json:"courses_name"`
	InstructorName string `json:"instructor_name"`
	Progress       string `json:"progress"`
	//NoOfStudents    int    `json:"no_of_students"`
	CoursesCategory string `json:"courses_category"`
}
