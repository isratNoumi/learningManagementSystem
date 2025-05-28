package models

import "time"

type InstructorResponse1 struct {
	CourseName     string `json:"coursename"`
	NoOfStudents   int    `json:"noofstudents"`
	InstructorName string `json:"instructor"`
}

type Instructor struct {
	ID        int       `json:"id"`
	UsersID   int       `json:"users_id"`
	CoursesID int       `json:"courses_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
