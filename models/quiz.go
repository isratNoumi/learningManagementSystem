package models

type QuizType struct {
	ID            int64  `json:"id"`
	Type          string `json:"type"`
	UnitDetailsID int64  `json:"unit_details_id"`
}

type QuizOptions struct {
	ID            int64  `json:"id"`
	UnitDetailsID int64  `json:"unit_details_id"`
	Options       string `json:"options"`
}
