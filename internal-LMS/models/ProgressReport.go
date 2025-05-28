package models

import "time"

type ProgressReport struct {
	ID               int64     `json:"id,omitempty"`
	UsersID          int64     `json:"user_id,omitempty"`
	CoursesID        int64     `json:"course_id,omitempty"`
	EnrollmentDate   time.Time `json:"enrollment_date" gorm:"autoCreateTime;<-:create"`
	OverallScore     int64     `json:"overall_score,omitempty"`
	CompletionStatus string    `json:"completion_status" gorm:"default:in_progress"`
}
