package models

type CourseContentResult struct {
	Coursename string `json:"course_name"`
	Modulename string `json:"module_name"`
	Lessonname string `json:"lesson_name"`
	Unittype   string `json:"unit_type"`
}
