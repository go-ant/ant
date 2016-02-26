package models

type Pagination struct {
	PerPage    uint32 `json:"per_page"`
	Page       uint32 `json:"page"`
	TotalPages uint32 `json:"total"`
}
