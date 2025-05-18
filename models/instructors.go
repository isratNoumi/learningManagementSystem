package models

type Instructor struct {
	CourseName     string `json:"coursename"`
	NoOfStudents   int    `json:"noofstudents"`
	InstructorName string `json:"instructor"`
}
