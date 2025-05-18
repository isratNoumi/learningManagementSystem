package models

// Pagination contains pagination metadata
type Pagination struct {
	CurrentPage  int `json:"currentPage"`
	PageSize     int `json:"pageSize"`
	TotalRecords int `json:"totalRecords"`
	TotalPages   int `json:"totalPages"`
}

// Response wraps the API response
type Response struct {
	Data       []CourseDTO `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Links      []Link      `json:"links"`
}
